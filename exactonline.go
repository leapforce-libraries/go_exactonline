package exactonline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	types "github.com/leapforce-libraries/go_types"
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
	division int
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
func NewExactOnline(division int, clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*ExactOnline, error) {
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

func (eo *ExactOnline) InitToken() error {
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

func (eo *ExactOnline) FindSubscriptionsForAccount(ac *Account) error {
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

	//fmt.Println("FindSubscriptionsForAccount:", len(ac.Subscriptions))
	return nil
}

// wait assures the maximum of 300(?) api calls per minute dictated by exactonline's rate-limit
func (eo *ExactOnline) Wait() error {
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

	return nil
}

// generic methods
//

func (eo *ExactOnline) ReadRateLimitHeaders(res *http.Response) {
	//fmt.Println("X-RateLimit-Minutely-Remaining", res.Header.Get("X-RateLimit-Minutely-Remaining"))
	//fmt.Println("X-RateLimit-Minutely-Reset", res.Header.Get("X-RateLimit-Minutely-Reset"))
	remaining, errRem := strconv.Atoi(res.Header.Get("X-RateLimit-Minutely-Remaining"))
	reset, errRes := strconv.ParseInt(res.Header.Get("X-RateLimit-Minutely-Reset"), 10, 64)
	if errRem == nil && errRes == nil {
		eo.XRateLimitMinutelyRemaining = remaining
		eo.XRateLimitMinutelyReset = reset
	}
}

func (eo *ExactOnline) PrintError(res *http.Response) error {
	fmt.Println("Status", res.Status)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println("errUnmarshal1")
		return err
	}

	ee := ExactOnlineError{}

	err = json.Unmarshal(b, &ee)
	if err != nil {
		//fmt.Println("errUnmarshal1")
		return err
	}

	//fmt.Println(ee.Err.Message.Value)
	message := fmt.Sprintf("Server returned statuscode %v, error:%s", res.StatusCode, ee.Err.Message.Value)
	return &types.ErrorString{message}
}

func (eo *ExactOnline) Get(url string, model interface{}) (string, error) {
	err := eo.Wait()
	if err != nil {
		return "", err
	}

	eo.RequestCount++

	response := Response{}
	res, err := eo.oAuth2.Get(url, &response)
	if err != nil {
		if res != nil {
			return "", eo.PrintError(res)
		} else {
			return "", err
		}

	}

	eo.ReadRateLimitHeaders(res)

	err = json.Unmarshal(response.Data.Results, &model)
	if err != nil {
		return "", err
	}

	return response.Data.Next, nil
}

func (eo *ExactOnline) PutValues(url string, values map[string]string) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return eo.Put(url, buf)
}

func (eo *ExactOnline) PutBytes(url string, b []byte) error {
	return eo.Put(url, bytes.NewBuffer(b))
}

func (eo *ExactOnline) Put(url string, buf *bytes.Buffer) error {
	eo.RequestCount++

	res, err := eo.oAuth2.Put(url, buf, nil)
	if err != nil {
		if res != nil {
			return eo.PrintError(res)
		} else {
			return err
		}
	}

	eo.ReadRateLimitHeaders(res)

	return nil
}

func (eo *ExactOnline) PostValues(url string, values map[string]string, model interface{}) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return eo.Post(url, buf, model)
}

func (eo *ExactOnline) PostBytes(url string, b []byte, model interface{}) error {
	return eo.Post(url, bytes.NewBuffer(b), model)
}

func (eo *ExactOnline) Post(url string, buf *bytes.Buffer, model interface{}) error {
	eo.RequestCount++

	response := ResponseSingle{}
	res, err := eo.oAuth2.Post(url, buf, &response)
	if err != nil {
		if res != nil {
			return eo.PrintError(res)
		} else {
			return err
		}
	}

	eo.ReadRateLimitHeaders(res)

	defer res.Body.Close()

	err = json.Unmarshal(response.Data, &model)
	if err != nil {
		return err
	}

	return nil
}

func (eo *ExactOnline) Delete(url string) error {
	eo.RequestCount++

	res, err := eo.oAuth2.Delete(url, nil, nil)
	if err != nil {
		if res != nil {
			return eo.PrintError(res)
		} else {
			return err
		}
	}

	eo.ReadRateLimitHeaders(res)

	return nil
}
