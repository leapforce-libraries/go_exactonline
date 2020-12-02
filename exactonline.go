package exactonline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"
	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName         string = "ExactOnline"
	apiURL          string = "https://start.exactonline.nl/api/v1"
	authURL         string = "https://start.exactonline.nl/api/oauth2/auth"
	tokenURL        string = "https://start.exactonline.nl/api/oauth2/token"
	tokenHttpMethod string = http.MethodPost
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// ExactOnline stores ExactOnline configuration
//
type ExactOnline struct {
	division int32
	oAuth2   *oauth2.OAuth2

	// data
	Contacts          []Contact
	Accounts          []Account
	SubscriptionTypes []SubscriptionType
	Subscriptions     []Subscription
	SubscriptionLines []SubscriptionLine
	Divisions         []Division
	Items             []Item
	//Token             *Token

	// rate limit
	XRateLimitMinutelyRemaining int
	XRateLimitMinutelyReset     int64
	RequestCount                int64
	//IsLive                      bool
}

// methods
//
func NewExactOnline(division int32, clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*ExactOnline, *errortools.Error) {
	eo := ExactOnline{}
	eo.division = division

	eo.RequestCount = 0

	config := oauth2.OAuth2Config{
		ApiName:         apiName,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHttpMethod,
	}
	eo.oAuth2 = oauth2.NewOAuth(config, bigQuery, isLive)
	return &eo, nil
}

func (eo *ExactOnline) baseURL() string {
	return fmt.Sprintf("%s/%v", apiURL, eo.division)
}

func (eo *ExactOnline) apiURL() string {
	return apiURL
}

func (eo *ExactOnline) InitToken() *errortools.Error {
	return eo.oAuth2.InitToken()
}

// Response represents highest level of exactonline api response
//
type Response struct {
	Data Results `json:"d"`
}

// ResponseSingle represents highest level of exactonline api response
//
type ResponseSingle struct {
	Data json.RawMessage `json:"d"`
}

// Results represents second highest level of exactonline api response
//
type Results struct {
	Results json.RawMessage `json:"results"`
	Next    string          `json:"__next"`
}

func (eo *ExactOnline) FindSubscriptionsForAccount(ac *Account) {
	for _, s := range eo.Subscriptions {
		if ac.ID == s.OrderedBy {
			for _, sl := range eo.SubscriptionLines {
				if s.EntryID == sl.EntryID {
					s.SubscriptionLines = append(s.SubscriptionLines, sl)
				}
			}

			ac.Subscriptions = append(ac.Subscriptions, s)
		}
	}
}

// wait assures the maximum of 300(?) api calls per minute dictated by exactonline's rate-limit
func (eo *ExactOnline) Wait() {
	if eo.XRateLimitMinutelyRemaining < 1 {
		reset := time.Unix(eo.XRateLimitMinutelyReset/1000, 0)
		ms := reset.Sub(time.Now()).Milliseconds()

		if ms > 0 {
			fmt.Println("eo.XRateLimitMinutelyReset:", eo.XRateLimitMinutelyReset)
			fmt.Println("reset:", reset)
			fmt.Println("waiting ms:", ms)
			time.Sleep(time.Duration(ms+1000) * time.Millisecond)
		}
	}
}

// generic methods
//

func (eo *ExactOnline) ReadRateLimitHeaders(res *http.Response) {
	remaining, errRem := strconv.Atoi(res.Header.Get("X-RateLimit-Minutely-Remaining"))
	reset, errRes := strconv.ParseInt(res.Header.Get("X-RateLimit-Minutely-Reset"), 10, 64)
	if errRem == nil && errRes == nil {
		eo.XRateLimitMinutelyRemaining = remaining
		eo.XRateLimitMinutelyReset = reset
	}
}

func (eo *ExactOnline) Get(url string, model interface{}) (string, *errortools.Error) {
	eo.Wait()

	eo.RequestCount++

	response := Response{}
	ee := ExactOnlineError{}
	_, res, e := eo.oAuth2.Get(url, &response, &ee)

	if e != nil {
		if ee.Err.Message.Value != "" {
			e.SetMessage(ee.Err.Message.Value)
		}

		return "", e
	}

	eo.ReadRateLimitHeaders(res)

	err := json.Unmarshal(response.Data.Results, &model)
	if err != nil {
		e.SetMessage(err)
		return "", e
	}

	return response.Data.Next, nil
}

func (eo *ExactOnline) PutValues(url string, values map[string]string) *errortools.Error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return eo.Put(url, buf)
}

func (eo *ExactOnline) PutBytes(url string, b []byte) *errortools.Error {
	return eo.Put(url, bytes.NewBuffer(b))
}

func (eo *ExactOnline) Put(url string, buf *bytes.Buffer) *errortools.Error {
	eo.RequestCount++

	ee := ExactOnlineError{}
	_, res, e := eo.oAuth2.Put(url, buf, nil, &ee)

	if e != nil {
		if ee.Err.Message.Value != "" {
			e.SetMessage(ee.Err.Message.Value)
		}

		return e
	}

	eo.ReadRateLimitHeaders(res)

	return nil
}

func (eo *ExactOnline) PostValues(url string, values map[string]string, model interface{}) *errortools.Error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return eo.Post(url, buf, model)
}

func (eo *ExactOnline) PostBytes(url string, b []byte, model interface{}) *errortools.Error {
	return eo.Post(url, bytes.NewBuffer(b), model)
}

func (eo *ExactOnline) Post(url string, buf *bytes.Buffer, model interface{}) *errortools.Error {
	eo.RequestCount++

	ee := ExactOnlineError{}
	response := ResponseSingle{}
	_, res, e := eo.oAuth2.Post(url, buf, &response, &ee)

	if e != nil {
		if ee.Err.Message.Value != "" {
			e.SetMessage(ee.Err.Message.Value)
		}

		return e
	}

	eo.ReadRateLimitHeaders(res)

	defer res.Body.Close()

	err := json.Unmarshal(response.Data, &model)
	if err != nil {
		e.SetMessage(err)
		return e
	}

	return nil
}

func (eo *ExactOnline) Delete(url string) *errortools.Error {
	eo.RequestCount++

	ee := ExactOnlineError{}
	_, res, e := eo.oAuth2.Delete(url, nil, nil, &ee)

	if e != nil {
		if ee.Err.Message.Value != "" {
			e.SetMessage(ee.Err.Message.Value)
		}

		return e
	}

	eo.ReadRateLimitHeaders(res)

	return nil
}
