package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// EmploymentInternalRate stores EmploymentInternalRate from exactonline
//
type EmploymentInternalRate struct {
	ID               types.GUID  `json:"ID"`
	Created          *types.Date `json:"Created,omitempty"`
	Creator          types.GUID  `json:"Creator"`
	CreatorFullName  string      `json:"CreatorFullName"`
	Division         int64       `json:"Division"`
	Employee         types.GUID  `json:"Employee"`
	EmployeeFullName string      `json:"EmployeeFullName"`
	EmployeeHID      int64       `json:"EmployeeHID"`
	Employment       types.GUID  `json:"Employment"`
	EmploymentHID    int64       `json:"EmploymentHID"`
	EndDate          *types.Date `json:"EndDate,omitempty"`
	InternalRate     float64     `json:"InternalRate"`
	Modified         *types.Date `json:"Modified,omitempty"`
	Modifier         types.GUID  `json:"Modifier"`
	ModifierFullName string      `json:"ModifierFullName"`
	StartDate        *types.Date `json:"StartDate,omitempty"`
}

func (eo *ExactOnline) GetEmploymentInternalRatesInternal(filter string) (*[]EmploymentInternalRate, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", EmploymentInternalRate{})
	urlStr := fmt.Sprintf("%s/project/EmploymentInternalRates?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	employmentInternalRates := []EmploymentInternalRate{}

	for urlStr != "" {
		ac := []EmploymentInternalRate{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetEmploymentInternalRatesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		employmentInternalRates = append(employmentInternalRates, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &employmentInternalRates, nil
}

func (eo *ExactOnline) GetEmploymentInternalRates() (*[]EmploymentInternalRate, *errortools.Error) {
	acc, err := eo.GetEmploymentInternalRatesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
