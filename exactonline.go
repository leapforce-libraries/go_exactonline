package exactonline

import (
	"bigquerytools"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	// get eContacts
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
	// get eContacts
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
	// get eAccounts
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

	// print
	//jsonString, _ = json.Marshal(eo.Accounts)
	//fmt.Println(string(jsonString))

	//
	// get eSubscriptions
	//
	errS := eo.getSubscriptions()
	if errS != nil {
		log.Fatal(errS)
	}
	fmt.Println("#eSubscriptions: ", len(eo.Subscriptions))

	b := make([]interface{}, len(eo.Subscriptions))
	for i := range eo.Subscriptions {
		b[i] = eo.Subscriptions[i]
	}
	/*
		bq := new(bigquery.BigQuery)
		SliceToBigQuery(b, eSubscription{})*/

	// print
	//jsonString, _ = json.Marshal(eo.Subscriptions)
	//fmt.Println(string(jsonString))

	//
	// get eSubscriptionLines
	//
	errSL := eo.getSubscriptionLines()
	if errSL != nil {
		log.Fatal(errSL)
	}
	fmt.Println("#eSubscriptionLines: ", len(eo.SubscriptionLines))

	// print
	//jsonString, _ = json.Marshal(eo.SubscriptionLines)
	//fmt.Println(string(jsonString))

	return nil
}

func (eo *ExactOnline) getHttpClient() (*http.Client, error) {
	//ctx := context.Background()

	err := eo.ValidateToken()
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

func (eo *ExactOnline) getSubscriptions() error {
	urlStr := fmt.Sprintf("https://start.exactonline.nl/api/v1/%s/subscription/Subscriptions", strconv.Itoa(eo.Me.CurrentDivision))

	_, err := eo.get(urlStr, &eo.Subscriptions)
	if err != nil {
		return err
	}

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

	// Check HTTP StatusCode
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", &types.ErrorString{fmt.Sprintf("Server returned statuscode %v", res.StatusCode)}
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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	//fmt.Println(res)

	return nil
}

func (eo *ExactOnline) post(url string, values map[string]string, model interface{}) error {
	client, errClient := eo.getHttpClient()
	if errClient != nil {
		return errClient
	}

	//jsonValue, _ := json.Marshal(values)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(values)

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		fmt.Println("errNewRequest")
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("errDo")
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))

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
