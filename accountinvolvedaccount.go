package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// AccountInvolvedAccount stores AccountInvolvedAccount from exactonline
//
type AccountInvolvedAccount struct {
	Account     types.GUID  `json:"Account"`
	AccountName string      `json:"AccountName"`
	Created     *types.Date `json:"Created,omitempty"`
	Creator     types.GUID  `json:"Creator"`
}

/*
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
}*/

func (eo *ExactOnline) GetAccountInvolvedAccountsInternal(filter string) (*[]AccountInvolvedAccount, error) {
	selectFields := GetJsonTaggedFieldNames(AccountInvolvedAccount{})
	urlStr := fmt.Sprintf("%s%s/accountancy/AccountInvolvedAccounts?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	accountInvolvedAccounts := []AccountInvolvedAccount{}

	for urlStr != "" {
		ac := []AccountInvolvedAccount{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetAccountInvolvedAccountsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		accountInvolvedAccounts = append(accountInvolvedAccounts, ac...)

		urlStr = str
	}

	return &accountInvolvedAccounts, nil
}

func (eo *ExactOnline) GetAccountInvolvedAccounts() (*[]AccountInvolvedAccount, error) {
	acc, err := eo.GetAccountInvolvedAccountsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
