package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// ActiveEmployment stores ActiveEmployment from exactonline
//
type ActiveEmployment struct {
	ID                       types.GUID  `json:"ID"`
	AverageDaysPerWeek       float64     `json:"AverageDaysPerWeek"`
	AverageHoursPerWeek      float64     `json:"AverageHoursPerWeek"`
	Contract                 types.GUID  `json:"Contract"`
	ContractDocument         types.GUID  `json:"ContractDocument"`
	ContractEndDate          *types.Date `json:"ContractEndDate,omitempty"`
	ContractProbationEndDate *types.Date `json:"ContractProbationEndDate,omitempty"`
	ContractProbationPeriod  int64       `json:"ContractProbationPeriod"`
	ContractStartDate        *types.Date `json:"ContractStartDate,omitempty"`
	ContractType             int64       `json:"ContractType"`
	ContractTypeDescription  string      `json:"ContractTypeDescription"`
	Created                  *types.Date `json:"Created,omitempty"`
	Creator                  types.GUID  `json:"Creator"`
	CreatorFullName          string      `json:"CreatorFullName"`
	Department               types.GUID  `json:"Department"`
	DepartmentCode           string      `json:"DepartmentCode"`
	DepartmentDescription    string      `json:"DepartmentDescription"`
	Division                 int64       `json:"Division"`
}

func (eo *ExactOnline) GetActiveEmploymentsInternal(filter string) (*[]ActiveEmployment, error) {
	selectFields := GetJsonTaggedFieldNames(ActiveEmployment{})
	urlStr := fmt.Sprintf("%s%s/payroll/ActiveEmployments?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	activeEmployments := []ActiveEmployment{}

	for urlStr != "" {
		ac := []ActiveEmployment{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetActiveEmploymentsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		activeEmployments = append(activeEmployments, ac...)

		urlStr = str

		urlStr = ""
	}

	return &activeEmployments, nil
}

func (eo *ExactOnline) GetActiveEmployments() (*[]ActiveEmployment, error) {
	acc, err := eo.GetActiveEmploymentsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
