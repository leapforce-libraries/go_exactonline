package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// AccountInvolvedAccount stores AccountInvolvedAccount from exactonline
//
type AccountInvolvedAccount struct {
	Account                                      types.GUID  `json:"Account"`
	AccountName                                  string      `json:"AccountName"`
	Created                                      *types.Date `json:"Created,omitempty"`
	Creator                                      types.GUID  `json:"Creator"`
	CreatorFullName                              string      `json:"CreatorFullName"`
	Division                                     int64       `json:"Division"`
	InvolvedAccount                              types.GUID  `json:"InvolvedAccount"`
	InvolvedAccountRelationTypeDescription       string      `json:"InvolvedAccountRelationTypeDescription"`
	InvolvedAccountRelationTypeDescriptionTermId int64       `json:"InvolvedAccountRelationTypeDescriptionTermId"`
	InvolvedAccountRelationTypeId                int64       `json:"InvolvedAccountRelationTypeId"`
	Modified                                     *types.Date `json:"Modified,omitempty"`
	Modifier                                     types.GUID  `json:"Modifier"`
	ModifierFullName                             string      `json:"ModifierFullName"`
	Notes                                        string      `json:"Notes"`
}

func (eo *ExactOnline) GetAccountInvolvedAccountsInternal(filter string) (*[]AccountInvolvedAccount, error) {
	selectFields := utilities.GetTaggedFieldNames("json", AccountInvolvedAccount{})
	urlStr := fmt.Sprintf("%s/accountancy/AccountInvolvedAccounts?$select=%s", eo.baseURL(), selectFields)
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
