package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// RevenueList stores RevenueList from exactonline
//
type RevenueList struct {
	Period int32   `json:"Period"`
	Year   int32   `json:"Year"`
	Amount float64 `json:"Amount"`
}

func (eo *ExactOnline) GetRevenueListsInternal(filter string) (*[]RevenueList, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", RevenueList{})
	urlStr := fmt.Sprintf("%s/read/financial/RevenueList?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	revenueLists := []RevenueList{}

	for urlStr != "" {
		its := []RevenueList{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetRevenueListsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		revenueLists = append(revenueLists, its...)

		urlStr = str
		//urlStr = ""
	}

	return &revenueLists, nil
}

func (eo *ExactOnline) GetRevenueLists() (*[]RevenueList, *errortools.Error) {
	acc, err := eo.GetRevenueListsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
