package exactonline

import (
	"fmt"
	"strconv"
	"strings"

	types "github.com/Leapforce-nl/go_types"
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
	Status                 string     `json:"Status"`
	AccountManager         types.GUID `json:"AccountManager"`
	AccountManagerFullName string     `json:"AccountManagerFullName"`
	MainContact            types.GUID `json:"MainContact"`
	Subscriptions          []Subscription
}

// AccountBigQuery stores account from exactonline
//
type AccountBigQuery struct {
	ID                     string
	Name                   string
	ChamberOfCommerce      string
	AddressLine1           string
	Postcode               string
	City                   string
	State                  string
	Country                string
	Status                 string
	AccountManager         string
	AccountManagerFullName string
	MainContact            string
}

// ToBigQuery convert Subscription to SubscriptionBq
//
func (a *Account) ToBigQuery() *AccountBigQuery {
	return &AccountBigQuery{
		a.ID.String(),
		a.Name,
		a.ChamberOfCommerce,
		a.AddressLine1,
		a.Postcode,
		a.City,
		a.State,
		a.Country,
		a.Status,
		a.AccountManager.String(),
		a.AccountManagerFullName,
		a.MainContact.String(),
	}
}

var oldAccount *Account

// SaveValues saves current values in local copy of Account
//
func (a *Account) SaveValues(inserted bool) {
	oldAccount = nil
	if !inserted {
		oldAccount = new(Account)
		oldAccount.Name = a.Name
		oldAccount.ChamberOfCommerce = a.ChamberOfCommerce
		oldAccount.AddressLine1 = a.AddressLine1
		oldAccount.Postcode = a.Postcode
		oldAccount.City = a.City
		oldAccount.State = a.State
		oldAccount.Country = a.Country
		oldAccount.Status = a.Status
		oldAccount.AccountManagerFullName = a.AccountManagerFullName
	}
}

// Values return comma separated values of Account
//
/*func (a *Account) Values() string {
	return fmt.Sprintf("Name: %s, ChamberOfCommerce: %s, AddressLine1: %s, Postcode: %s, c: %s, State: %s, Country: %s",
		a.Name,
		a.ChamberOfCommerce,
		a.AddressLine1,
		a.Postcode,
		a.Postcode,
		a.State,
		a.Country)
}*/

func (a *Account) Values() (string, string) {
	old := ""
	new := ""

	if oldAccount == nil {
		new += ",Name:" + a.Name
	} else if oldAccount.Name != a.Name {
		old += ",Name:" + oldAccount.Name
		new += ",Name:" + a.Name
	}

	if oldAccount != nil {
		oldAccount.Name = strings.Trim(oldAccount.Name, " ")
		oldAccount.AddressLine1 = strings.Trim(oldAccount.AddressLine1, " ")
		oldAccount.Postcode = strings.Trim(oldAccount.Postcode, " ")
		oldAccount.City = strings.Trim(oldAccount.City, " ")
		oldAccount.State = strings.Trim(oldAccount.State, " ")
		oldAccount.Country = strings.Trim(oldAccount.Country, " ")
		oldAccount.Status = strings.Trim(oldAccount.Status, " ")
	}

	a.Name = strings.Trim(a.Name, " ")
	a.AddressLine1 = strings.Trim(a.AddressLine1, " ")
	a.Postcode = strings.Trim(a.Postcode, " ")
	a.City = strings.Trim(a.City, " ")
	a.State = strings.Trim(a.State, " ")
	a.Country = strings.Trim(a.Country, " ")
	a.Status = strings.Trim(a.Status, " ")

	if oldAccount == nil {
		new += ",ChamberOfCommerce:" + a.ChamberOfCommerce
	} else if oldAccount.ChamberOfCommerce != a.ChamberOfCommerce {
		old += ",ChamberOfCommerce:" + oldAccount.ChamberOfCommerce
		new += ",ChamberOfCommerce:" + a.ChamberOfCommerce
	}

	if oldAccount == nil {
		new += ",AddressLine1:" + a.AddressLine1
	} else if oldAccount.AddressLine1 != a.AddressLine1 {
		old += ",AddressLine1:" + oldAccount.AddressLine1
		new += ",AddressLine1:" + a.AddressLine1
	}

	if oldAccount == nil {
		new += ",Postcode:" + a.Postcode
	} else if oldAccount.Postcode != a.Postcode {
		old += ",Postcode:" + oldAccount.Postcode
		new += ",Postcode:" + a.Postcode
	}

	if oldAccount == nil {
		new += ",City:" + a.City
	} else if oldAccount.City != a.City {
		old += ",City:" + oldAccount.City
		new += ",City:" + a.City
	}

	if oldAccount == nil {
		new += ",State:" + a.State
	} else if oldAccount.State != a.State {
		old += ",State:" + oldAccount.State
		new += ",State:" + a.State
	}

	if oldAccount == nil {
		new += ",Country:" + a.Country
	} else if oldAccount.Country != a.Country {
		old += ",Country:" + oldAccount.Country
		new += ",Country:" + a.Country
	}

	if oldAccount == nil {
		new += ",Status:" + a.Status
	} else if oldAccount.Status != a.Status {
		old += ",Status:" + oldAccount.Status
		new += ",Status:" + a.Status
	}

	if oldAccount == nil {
		new += ",AccountManagerFullName:" + a.AccountManagerFullName
	} else if oldAccount.AccountManagerFullName != a.AccountManagerFullName {
		old += ",AccountManagerFullName:" + oldAccount.AccountManagerFullName
		new += ",AccountManagerFullName:" + a.AccountManagerFullName
	}

	old = strings.TrimLeft(old, ",")
	new = strings.TrimLeft(new, ",")

	return old, new
}

func (eo *ExactOnline) GetAccountsInternal(filter string) (*[]Account, error) {
	selectFields := GetJsonTaggedFieldNames(Account{})
	urlStr := fmt.Sprintf("%s%s/crm/Accounts?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	accounts := []Account{}

	for urlStr != "" {
		ac := []Account{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetAccountsInternal:", err)
			fmt.Println("url:", urlStr)
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
	urlStr := fmt.Sprintf("%s%s/crm/Accounts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), a.ID.String())

	data := make(map[string]string)
	data["Name"] = a.Name
	data["AddressLine1"] = a.AddressLine1
	data["Postcode"] = a.Postcode
	data["City"] = a.City
	data["State"] = "" //a.State
	data["Country"] = a.Country
	data["Status"] = a.Status
	//data["State"] = "ZH"   //a.StateName + "_updated3"
	//data["Country"] = "NL" //a.CountryName + "_updated3"
	//data["AccountManagerFullName"] = a.AccountManagerFullName + "_updated3"

	//fmt.Println("ID")
	//fmt.Println("Updated:", a.ID.String(), data["AddressLine1"])

	err := eo.Put(urlStr, data)
	if err != nil {
		fmt.Println("ERROR in UpdateAccount:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", data)
		return err
	}

	fmt.Println("\nUPDATED Account")
	fmt.Println("url:", urlStr)
	fmt.Println("data:", data)

	return nil
}

func (eo *ExactOnline) UpdateAccountMainContact(a *Account) error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), a.ID.String())

	data := make(map[string]string)
	data["MainContact"] = a.MainContact.String()

	err := eo.Put(urlStr, data)
	if err != nil {
		fmt.Println("ERROR in UpdateAccountMainContact:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", data)
		return err
	}

	fmt.Println("\nUPDATED Account (MainContact)")
	fmt.Println("url:", urlStr)
	fmt.Println("data:", data)

	return nil
}

func (eo *ExactOnline) InsertAccount(a *Account) error {
	urlStr := fmt.Sprintf("%s%s/crm/Accounts", eo.ApiUrl, strconv.Itoa(eo.Division))

	data := make(map[string]string)
	data["Name"] = a.Name
	data["AddressLine1"] = a.AddressLine1
	data["Postcode"] = a.Postcode
	data["City"] = a.City
	data["State"] = "" //a.State
	data["Country"] = a.Country
	data["Status"] = a.Status
	//data["State"] = a.StateName
	//data["Country"] = a.CountryName
	//data["AccountManagerFullName"] = a.AccountManagerFullName
	data["ChamberOfCommerce"] = a.ChamberOfCommerce

	//fmt.Println(urlStr)
	ac := Account{}

	err := eo.Post(urlStr, data, &ac)
	if err != nil {
		fmt.Println("ERROR in InsertAccount:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", data)
		return err
	}

	fmt.Println("\nINSERTED Account", ac.ID)
	fmt.Println("url:", urlStr)
	fmt.Println("data:", data)

	a.ID = ac.ID

	//fmt.Println("Inserted:", a.ID.String())

	//time.Sleep(1 * time.Second)

	return nil
}
