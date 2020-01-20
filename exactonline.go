package exactonline

import (
	"bigquerytools"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"types"
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
	Subscriptions     []Subscription
	SubscriptionLines []SubscriptionLine
	Divisions         []Division
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

//
// getAll retrieves all tables
//
func (eo *ExactOnline) GetAll() error {
	//
	// get eMe
	//
	errMe := eo.GetMe()
	if errMe != nil {
		log.Fatal(errMe)
	}

	// print
	fmt.Printf("CurrentDivision:")
	fmt.Println(eo.Me.CurrentDivision)

	//
	// get SubscriptionLines
	//
	errSL := eo.getSubscriptionLines()
	if errSL != nil {
		log.Fatal(errSL)
	}
	fmt.Println("#eSubscriptionLines: ", len(eo.SubscriptionLines))

	//
	// get Subscriptions
	//
	errS := eo.getSubscriptions()
	if errS != nil {
		log.Fatal(errS)
	}
	fmt.Println("#eSubscriptions: ", len(eo.Subscriptions))

	bs := make([]interface{}, len(eo.Subscriptions))
	for i := range eo.Subscriptions {
		bs[i] = eo.Subscriptions[i].ToBq()
	}

	errInsert2 := eo.BigQuery.InsertSlice(eo.BigQueryDataset, bs, SubscriptionBq{}, "subscriptions_temp")
	if errInsert2 != nil {
		log.Fatal(errInsert2)
	}

	//
	// get Divisions
	//

	errD := eo.getDivisions()
	if errD != nil {
		log.Fatal(errD)
	}
	/*for _, co := range eo.Contacts {
		//jsonString, _ := json.Marshal(a)
		//fmt.Println(string(jsonString))
		fmt.Println("Account:", co.Account.String(), "Contact:", co.ID.String())
	}*/
	fmt.Println("#eDivisions: ", len(eo.Divisions))

	//
	// get Contacts
	//
	errC := eo.getContacts()
	if errC != nil {
		log.Fatal(errC)
	}
	/*for _, co := range eo.Contacts {
		//jsonString, _ := json.Marshal(a)
		//fmt.Println(string(jsonString))
		oldValues := fmt.Sprintln(co)
		fmt.Println(oldValues)
		//fmt.Println("Account:", co.Account.String(), "Contact:", co.ID.String())
	}*/
	fmt.Println("#eContacts: ", len(eo.Contacts))

	// print
	//jsonString, _ := json.Marshal(eo.Contacts)
	//fmt.Println(string(jsonString))

	//
	// get Accounts
	//
	errA := eo.getAccounts()
	if errA != nil {
		log.Fatal(errA)
	}
	/*for _, a := range eo.Accounts {
		//jsonString, _ := json.Marshal(a)
		//fmt.Println(string(jsonString))
		fmt.Println(a.ID.String())
	}*/
	fmt.Println("#eAccounts: ", len(eo.Accounts))

	b1 := make([]interface{}, len(eo.Accounts))
	for i := range eo.Accounts {
		b1[i] = eo.Accounts[i].ToBq()
	}

	errInsert := eo.BigQuery.InsertSlice(eo.BigQueryDataset, b1, AccountBq{}, "accounts_temp")
	if errInsert != nil {
		log.Fatal(errInsert)
	}

	eo.GetSubscriptionsForAccounts()

	return nil
}

func (eo *ExactOnline) GetSubscriptionsForAccounts() {
	count := 0
	for _, a := range eo.Accounts {
		for _, s := range eo.Subscriptions {
			if a.ID == s.OrderedBy {
				a.Subscriptions = append(a.Subscriptions, s)
				count++
			}
		}
	}

	//fmt.Println("GetSubscriptionsForAccounts:", count)
}

// wait assures the maximum of 300(?) api calls per minute dictated by exactonline's rate-limit
func (eo *ExactOnline) wait() error {
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

func (eo *ExactOnline) getHttpClient() (*http.Client, error) {
	err := eo.wait()
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

	_, err := eo.get(urlStr, &me)
	if err != nil {
		return err
	}

	eo.Me = me[0]

	return nil
}

func (eo *ExactOnline) getSubscriptionLines() error {
	urlStr := fmt.Sprintf("https://start.exactonline.nl/api/v1/%s/subscription/SubscriptionLines", strconv.Itoa(eo.Me.CurrentDivision))

	_, err := eo.get(urlStr, &eo.SubscriptionLines)
	if err != nil {
		return err
	}

	return nil
}

//
// generic methods
//
func (eo *ExactOnline) readRateLimitHeaders(res *http.Response) {
	//fmt.Println("X-RateLimit-Minutely-Remaining", res.Header.Get("X-RateLimit-Minutely-Remaining"))
	//fmt.Println("X-RateLimit-Minutely-Reset", res.Header.Get("X-RateLimit-Minutely-Reset"))
	remaining, errRem := strconv.Atoi(res.Header.Get("X-RateLimit-Minutely-Remaining"))
	reset, errRes := strconv.ParseInt(res.Header.Get("X-RateLimit-Minutely-Reset"), 10, 64)
	if errRem == nil && errRes == nil {
		eo.XRateLimitMinutelyRemaining = remaining
		eo.XRateLimitMinutelyReset = reset
	}
}

func (eo *ExactOnline) get(url string, model interface{}) (string, error) {
	client, errClient := eo.getHttpClient()
	if errClient != nil {
		return "", errClient
	}

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

	eo.readRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("Status", res.Status)
		fmt.Println(url)
		fmt.Println(eo.Token.AccessToken)
		return "", &types.ErrorString{fmt.Sprintf("Server returned statuscode %v: %s", res.StatusCode, err.Error())}
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

func (eo *ExactOnline) put(url string, values map[string]string) error {
	client, errClient := eo.getHttpClient()
	if errClient != nil {
		return errClient
	}

	//jsonValue, _ := json.Marshal(values)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

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

	eo.readRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("Status", res.Status)
		fmt.Println(url, values)
		fmt.Println(eo.Token.AccessToken)
		return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v: %s", res.StatusCode, err.Error())}
	}

	//fmt.Println(res)

	return nil
}

func (eo *ExactOnline) post(url string, values map[string]string, model interface{}) error {
	client, errClient := eo.getHttpClient()
	if errClient != nil {
		return errClient
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

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

	eo.readRateLimitHeaders(res)

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("Status", res.Status)
		fmt.Println(url, values)
		fmt.Println(eo.Token.AccessToken)
		return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v: %s", res.StatusCode, err.Error())}
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := ResponseSingle{}

	errr := json.Unmarshal(b, &response)
	if errr != nil {
		fmt.Println("errUnmarshal1")
		return err
	}

	errrr := json.Unmarshal(response.Data, &model)
	if errrr != nil {
		fmt.Println("errUnmarshal2")
		return errrr
	}

	return nil
}
