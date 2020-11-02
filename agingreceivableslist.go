package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// AgingReceivablesList stores AgingReceivablesList from exactonline
//
type AgingReceivablesList struct {
	AccountId            types.GUID `json:"AccountId"`
	AccountCode          string     `json:"AccountCode"`
	AccountName          string     `json:"AccountName"`
	AgeGroup1            int32      `json:"AgeGroup1"`
	AgeGroup1Amount      float64    `json:"AgeGroup1Amount"`
	AgeGroup1Description string     `json:"AgeGroup1Description"`
	AgeGroup2            int32      `json:"AgeGroup2"`
	AgeGroup2Amount      float64    `json:"AgeGroup2Amount"`
	AgeGroup2Description string     `json:"AgeGroup2Description"`
	AgeGroup3            int32      `json:"AgeGroup3"`
	AgeGroup3Amount      float64    `json:"AgeGroup3Amount"`
	AgeGroup3Description string     `json:"AgeGroup3Description"`
	AgeGroup4            int32      `json:"AgeGroup4"`
	AgeGroup4Amount      float64    `json:"AgeGroup4Amount"`
	AgeGroup4Description string     `json:"AgeGroup4Description"`
	CurrencyCode         string     `json:"CurrencyCode"`
	TotalAmount          float64    `json:"TotalAmount"`
}

func (eo *ExactOnline) GetAgingReceivablesListsInternal(filter string) (*[]AgingReceivablesList, error) {
	selectFields := utilities.GetTaggedFieldNames("json", AgingReceivablesList{})
	urlStr := fmt.Sprintf("%s/read/financial/AgingReceivablesList?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	agingReceivablesLists := []AgingReceivablesList{}

	for urlStr != "" {
		ac := []AgingReceivablesList{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetAgingReceivablesListsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		agingReceivablesLists = append(agingReceivablesLists, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &agingReceivablesLists, nil
}

func (eo *ExactOnline) GetAgingReceivablesLists() (*[]AgingReceivablesList, error) {
	acc, err := eo.GetAgingReceivablesListsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
