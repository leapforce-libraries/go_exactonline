package exactonline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	types "github.com/Leapforce-nl/go_types"
)

// ExactOnline stores exactonline configuration
//
type ExactOnline struct {
	// config
	ClientID        string
	ClientSecret    string
	RedirectURL     string
	AuthURL         string
	TokenURL        string
	RefreshTokenKey string
	ApiUrl          string
	// bigquery
	BigQuery          *bigquerytools.BigQuery
	BigQueryDataset   string
	BigQueryTablename string
	InitCallback      callbackFunction
	// data
	Me                Me
	Contacts          []Contact
	Accounts          []Account
	SubscriptionTypes []SubscriptionType
	Subscriptions     []Subscription
	SubscriptionLines []SubscriptionLine
	Divisions         []Division
	Items             []Item
	Token             *Token
	// timer
	//LastApiCall time.Time
	//TimestampsTimestamps []time.Time
	// rate limit
	XRateLimitMinutelyRemaining int
	XRateLimitMinutelyReset     int64
}

type callbackFunction func()

// methods
//
func (eo *ExactOnline) Init() error {
	if eo.ApiUrl == "" {
		return &types.ErrorString{"ExactOnline ApiUrl not provided"}
	}
	/*if eo.Token == nil {
		return &errorString{"ExactOnline Token not provided"}
	}*/

	if !strings.HasSuffix(eo.ApiUrl, "/") {
		eo.ApiUrl = eo.ApiUrl + "/"
	}

	return nil
}

// GetJsonTaggedFieldNames returns comma separated string of
// fieldnames of struct having a json tag
//
func GetJsonTaggedFieldNames(model interface{}) string {
	val := reflect.ValueOf(model)
	list := ""
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			list += "," + field.Name
		}
	}

	list = strings.Trim(list, ",")

	return list
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
			fmt.Println("len(sub.SubscriptionLines)", len(s.SubscriptionLines))

			ac.Subscriptions = append(ac.Subscriptions, s)
		}
	}

	//fmt.Println("FindSubscriptionsForAccount:", len(ac.Subscriptions))
	return nil
}

// wait assures the maximum of 300(?) api calls per minute dictated by exactonline's rate-limit
func (eo *ExactOnline) Wait() error {
	if eo.XRateLimitMinutelyRemaining < 1 {
		reset := time.Unix(eo.XRateLimitMinutelyReset, 0)
		ms := reset.Sub(time.Now()).Milliseconds()

		if ms > 0 {
			fmt.Println("waiting ms:", ms)
			time.Sleep(time.Duration(ms+1000) * time.Millisecond)
		}
	}

	return nil

	/*

		maxCallsPerMinute := 60
		msPerMinute := int64(60500) // 60000 ms go in a minute, plus a small margin...
		len := len(eo.Timestamps)

		if len >= maxCallsPerMinute {
			ts := eo.Timestamps[len-maxCallsPerMinute]
			ms := time.Now().Sub(ts).Milliseconds()

			//fmt.Println(len, ms)

			if ms < msPerMinute {
				fmt.Println("waiting: ", (msPerMinute - ms), "ms")
				time.Sleep(time.Duration(msPerMinute-ms) * time.Millisecond)
			}
		}

		// add new timestamp
		eo.Timestamps = append(eo.Timestamps, time.Now())

		return nil*/
}

func (eo *ExactOnline) GetHttpClient() (*http.Client, error) {
	err := eo.Wait()
	if err != nil {
		return nil, err
	}

	err = eo.ValidateToken()
	if err != nil {
		return nil, err
	}

	return new(http.Client), nil
}

func (eo *ExactOnline) GetMe() error {
	urlStr := "https://start.exactonline.nl/api/v1/current/Me"

	me := []Me{}

	_, err := eo.Get(urlStr, &me)
	if err != nil {
		return err
	}

	eo.Me = me[0]

	return nil
}

//
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

func (eo *ExactOnline) Get(url string, model interface{}) (string, error) {
	client, errClient := eo.GetHttpClient()
	if errClient != nil {
		return "", errClient
	}

	//fmt.Println(url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	// Add authorization token to header
	var bearer = "Bearer " + eo.Token.AccessToken
	req.Header.Add("authorization", bearer)
	req.Header.Set("Accept", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	eo.ReadRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("ERROR in Get")
		fmt.Println(url)
		fmt.Println("StatusCode", res.StatusCode)
		fmt.Println(eo.Token.AccessToken)
		return "", eo.PrintError(res)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	response := Response{}

	errr := json.Unmarshal(b, &response)
	if errr != nil {
		return "", err
	}

	errrr := json.Unmarshal(response.Data.Results, &model)
	if errrr != nil {
		return "", errrr
	}

	return response.Data.Next, nil
}

func (eo *ExactOnline) Put(url string, values map[string]string) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return eo.PutBuffer(url, buf)
}

func (eo *ExactOnline) PutBytes(url string, b []byte) error {
	return eo.PutBuffer(url, bytes.NewBuffer(b))
}

func (eo *ExactOnline) PrintError(res *http.Response) error {
	fmt.Println("Status", res.Status)

	b, err := ioutil.ReadAll(res.Body)
	return err

	ee := ExactOnlineError{}

	err = json.Unmarshal(b, &ee)
	if err != nil {
		fmt.Println("errUnmarshal1")
		//errortools.Fatal(err)
	}

	fmt.Println(ee.Err.Message.Value)

	return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v", res.StatusCode)}
}

func (eo *ExactOnline) PutBuffer(url string, buf *bytes.Buffer) error {
	client, errClient := eo.GetHttpClient()
	if errClient != nil {
		return errClient
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}
	// Add authorization token to header
	var bearer = "Bearer " + eo.Token.AccessToken
	req.Header.Add("authorization", bearer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	eo.ReadRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("ERROR in Put")
		fmt.Println(url)
		fmt.Println("StatusCode", res.StatusCode)
		fmt.Println(eo.Token.AccessToken)
		return eo.PrintError(res)
	}

	//fmt.Println(res)

	return nil
}

func (eo *ExactOnline) Post(url string, values map[string]string, model interface{}) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	return eo.PostBuffer(url, buf, model)
}

func (eo *ExactOnline) PostBytes(url string, b []byte, model interface{}) error {
	return eo.PostBuffer(url, bytes.NewBuffer(b), model)
}

func (eo *ExactOnline) PostBuffer(url string, buf *bytes.Buffer, model interface{}) error {
	client, errClient := eo.GetHttpClient()
	if errClient != nil {
		return errClient
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		fmt.Println("errNewRequest")
		return err
	}
	// Add authorization token to header
	var bearer = "Bearer " + eo.Token.AccessToken
	req.Header.Add("authorization", bearer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("errDo")
		return err
	}

	eo.ReadRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("ERROR in Post")
		fmt.Println(url)
		fmt.Println("StatusCode", res.StatusCode)
		fmt.Println(eo.Token.AccessToken)
		return eo.PrintError(res)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := ResponseSingle{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		fmt.Println("errUnmarshal1")
		return eo.PrintError(res)
	}

	err = json.Unmarshal(response.Data, &model)
	if err != nil {
		fmt.Println("errUnmarshal2")
		return err
	}

	return nil
}

func (eo *ExactOnline) Delete(url string) error {
	client, err := eo.GetHttpClient()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println("errNewRequest")
		return err
	}
	// Add authorization token to header
	var bearer = "Bearer " + eo.Token.AccessToken
	req.Header.Add("authorization", bearer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("errDo")
		return err
	}

	eo.ReadRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("ERROR in Delete")
		fmt.Println(url)
		fmt.Println("StatusCode", res.StatusCode)
		fmt.Println(eo.Token.AccessToken)
		return err
	}

	defer res.Body.Close()

	return nil
}
