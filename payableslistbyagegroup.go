package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
)

// PayablesListByAgeGroup stores PayablesListByAgeGroup from exactonline
//
type PayablesListByAgeGroup struct {
	HID                string      `json:"HID "`
	AccountCode        string      `json:"AccountCode"`
	AccountId          types.GUID  `json:"AccountId"`
	AccountName        string      `json:"AccountName"`
	Amount             float64     `json:"Amount"`
	AmountInTransit    float64     `json:"AmountInTransit"`
	ApprovalStatus     int16       `json:"ApprovalStatus"`
	CurrencyCode       string      `json:"CurrencyCode"`
	Description        string      `json:"Description"`
	DueDate            *types.Date `json:"DueDate"`
	EntryNumber        int32       `json:"EntryNumber"`
	Id                 types.GUID  `json:"Id"`
	InvoiceDate        *types.Date `json:"InvoiceDate"`
	InvoiceNumber      int32       `json:"InvoiceNumber"`
	JournalCode        string      `json:"JournalCode"`
	JournalDescription string      `json:"JournalDescription"`
	YourRef            string      `json:"YourRef"`
}

func (eo *ExactOnline) GetPayablesListByAgeGroupsInternal(ageGroup int, filter string) (*[]PayablesListByAgeGroup, error) {
	selectFields := GetJsonTaggedFieldNames(PayablesListByAgeGroup{})
	urlStr := fmt.Sprintf("%s%v/read/financial/PayablesListByAgeGroup?ageGroup=%v&$select=%s", eo.ApiUrl, eo.Division, ageGroup, selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	payablesListByAgeGroups := []PayablesListByAgeGroup{}

	for urlStr != "" {
		ac := []PayablesListByAgeGroup{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetPayablesListByAgeGroupsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		payablesListByAgeGroups = append(payablesListByAgeGroups, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &payablesListByAgeGroups, nil
}

func (eo *ExactOnline) GetPayablesListByAgeGroups(ageGroup int) (*[]PayablesListByAgeGroup, error) {
	acc, err := eo.GetPayablesListByAgeGroupsInternal(ageGroup, "")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
