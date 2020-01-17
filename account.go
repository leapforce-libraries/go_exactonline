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
	StateName              string     `json:"StateName"`
	Country                string     `json:"Country"`
	CountryName            string     `json:"CountryName"`
	AccountManager         types.GUID `json:"AccountManager,omitempty"`
	AccountManagerFullName string     `json:"AccountManagerFullName"`
	MainContact            types.GUID `json:"MainContact"`
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

func (eo *ExactOnline) getAccounts() error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision))

	eo.Accounts = []Account{}

	for urlStr != "" {
		ac := []Account{}

		str, err := eo.get(urlStr, &ac)
		if err != nil {
			return err
		}

		eo.Accounts = append(eo.Accounts, ac...)
		//fmt.Println(len(eo.Accounts))

		urlStr = str
		//urlStr = "" //temp
	}

	return nil
}

func (eo *ExactOnline) UpdateAccount(a *Account) error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), a.ID.String())

	data := make(map[string]string)
	data["Name"] = a.Name
	data["AddressLine1"] = a.AddressLine1
	data["Postcode"] = a.Postcode
	data["City"] = a.City
	//data["State"] = "ZH"   //a.StateName + "_updated3"
	//data["Country"] = "NL" //a.CountryName + "_updated3"
	//data["AccountManagerFullName"] = a.AccountManagerFullName + "_updated3"

	//fmt.Println("ID")
	//fmt.Println("Updated:", a.ID.String(), data["AddressLine1"])

	//fmt.Println(urlStr)

	err := eo.put(urlStr, data)
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

	err := eo.put(urlStr, data)
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
	//data["State"] = a.StateName
	//data["Country"] = a.CountryName
	//data["AccountManagerFullName"] = a.AccountManagerFullName
	data["ChamberOfCommerce"] = a.ChamberOfCommerce

	//fmt.Println(urlStr)
	ac := Account{}

	err := eo.post(urlStr, data, &ac)
	if err != nil {
		fmt.Println(err)
		return err
	}

	a.ID = ac.ID

	//fmt.Println("Inserted:", a.ID.String())

	//time.Sleep(1 * time.Second)

	return nil
}
