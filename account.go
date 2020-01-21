package exactonline

import (
	"fmt"
	"strconv"
	"types"
)

// Account stores account from exactonline
//
type Account struct {
	ID                     types.GUID `json:"ID"`
	Name                   string     `json:"Name"`
	ChamberOfCommerce      string     `json:"ChamberOfCommerce"`
	AddressLine1           string     `json:"AddressLine1"`
	Postcode               string     `json:"Postcode"`
	City                   string     `json:"City"`
	State                  string     `json:"State"`
	Country                string     `json:"Country"`
	AccountManager         types.GUID `json:"AccountManager"`
	AccountManagerFullName string     `json:"AccountManagerFullName"`
	MainContact            types.GUID `json:"MainContact"`
	Subscriptions          []Subscription
}

// AccountBq stores account from exactonline
//
type AccountBq struct {
	ID                     string
	Name                   string
	ChamberOfCommerce      string
	AddressLine1           string
	Postcode               string
	City                   string
	State                  string
	Country                string
	AccountManager         string
	AccountManagerFullName string
	MainContact            string
}

// ToBq convert Subscription to SubscriptionBq
//
func (a *Account) ToBq() *AccountBq {
	return &AccountBq{
		a.ID.String(),
		a.Name,
		a.ChamberOfCommerce,
		a.AddressLine1,
		a.Postcode,
		a.City,
		a.State,
		a.Country,
		a.AccountManager.String(),
		a.AccountManagerFullName,
		a.MainContact.String(),
	}
}

// Values return comma separated values of Account
//
func (a *Account) Values() string {
	return fmt.Sprintf("Name: %s, ChamberOfCommerce: %s, AddressLine1: %s, Postcode: %s, c: %s, State: %s, Country: %s",
		a.Name,
		a.ChamberOfCommerce,
		a.AddressLine1,
		a.Postcode,
		a.Postcode,
		a.State,
		a.Country)
}

func (eo *ExactOnline) GetAccountsInternal(filter string) (*[]Account, error) {
	selectFields := GetJsonTaggedFieldNames(Account{})
	urlStr := fmt.Sprintf("%s%s/crm/Accounts?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	accounts := []Account{}

	for urlStr != "" {
		ac := []Account{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, ac...)

		urlStr = str
	}

	return &accounts, nil
}

func (eo *ExactOnline) GetAccounts() error {
	acc, err := eo.GetAccountsInternal("")
	if err != nil {
		return err
	}
	eo.Accounts = *acc

	return nil
}

func (eo ExactOnline) GetAccountsByChamberOfCommerce(chamberOfCommerce string) ([]Account, error) {
	filter := fmt.Sprintf("ChamberOfCommerce eq '%s'", chamberOfCommerce)
	accounts := []Account{}

	acc, err := eo.GetAccountsInternal(filter)
	if err != nil {
		return accounts, err
	}
	accounts = *acc

	return accounts, nil
}

func (eo *ExactOnline) UpdateAccount(a *Account) error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), a.ID.String())

	data := make(map[string]string)
	data["Name"] = a.Name
	data["AddressLine1"] = a.AddressLine1
	data["Postcode"] = a.Postcode
	data["City"] = a.City
	data["Country"] = a.Country
	//data["State"] = "ZH"   //a.StateName + "_updated3"
	//data["Country"] = "NL" //a.CountryName + "_updated3"
	//data["AccountManagerFullName"] = a.AccountManagerFullName + "_updated3"

	//fmt.Println("ID")
	//fmt.Println("Updated:", a.ID.String(), data["AddressLine1"])

	fmt.Println("update", urlStr, a.Country, a.Name)

	err := eo.Put(urlStr, data)
	if err != nil {
		return err
	}

	//time.Sleep(1 * time.Second)

	return nil
}

func (eo *ExactOnline) UpdateAccountMainContact(a *Account) error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), a.ID.String())

	data := make(map[string]string)
	data["MainContact"] = a.MainContact.String()

	err := eo.Put(urlStr, data)
	if err != nil {
		return err
	}

	return nil
}

func (eo *ExactOnline) InsertAccount(a *Account) error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision))

	data := make(map[string]string)
	data["Name"] = a.Name
	data["AddressLine1"] = a.AddressLine1
	data["Postcode"] = a.Postcode
	data["City"] = a.City
	data["Country"] = a.Country
	//data["State"] = a.StateName
	//data["Country"] = a.CountryName
	//data["AccountManagerFullName"] = a.AccountManagerFullName
	data["ChamberOfCommerce"] = a.ChamberOfCommerce

	//fmt.Println(urlStr)
	ac := Account{}

	fmt.Println("insert", urlStr, a.Country, a.Name)

	err := eo.Post(urlStr, data, &ac)
	if err != nil {
		fmt.Println(err)
		return err
	}

	a.ID = ac.ID

	//fmt.Println("Inserted:", a.ID.String())

	//time.Sleep(1 * time.Second)

	return nil
}
