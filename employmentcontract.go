package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// EmploymentContract stores EmploymentContract from exactonline
//
type EmploymentContract struct {
	ID                           types.GUID  `json:"ID"`
	ContractFlexPhase            int32       `json:"ContractFlexPhase"`
	ContractFlexPhaseDescription string      `json:"ContractFlexPhaseDescription"`
	Created                      *types.Date `json:"Created"`
	Creator                      types.GUID  `json:"Creator"`
	CreatorFullName              string      `json:"CreatorFullName"`
	Division                     int32       `json:"Division"`
	Document                     types.GUID  `json:"Document"`
	Employee                     types.GUID  `json:"Employee"`
	EmployeeFullName             string      `json:"EmployeeFullName"`
	EmployeeHID                  int32       `json:"EmployeeHID"`
	EmployeeType                 int32       `json:"EmployeeType"`
	EmployeeTypeDescription      string      `json:"EmployeeTypeDescription"`
	Employment                   types.GUID  `json:"Employment"`
	EmploymentHID                int32       `json:"EmploymentHID"`
	EndDate                      *types.Date `json:"EndDate"`
	Modified                     *types.Date `json:"Modified"`
	Modifier                     types.GUID  `json:"Modifier"`
	ModifierFullName             string      `json:"ModifierFullName"`
	Notes                        string      `json:"Notes"`
	ProbationEndDate             *types.Date `json:"ProbationEndDate"`
	ProbationPeriod              int32       `json:"ProbationPeriod"`
	ReasonContract               int32       `json:"ReasonContract"`
	ReasonContractDescription    string      `json:"ReasonContractDescription"`
	Sequence                     int32       `json:"Sequence"`
	StartDate                    *types.Date `json:"StartDate"`
	Type                         int32       `json:"Type"`
	TypeDescription              string      `json:"TypeDescription"`
}

func (eo *ExactOnline) GetEmploymentContractsInternal(filter string) (*[]EmploymentContract, error) {
	selectFields := utilities.GetTaggedFieldNames("json", EmploymentContract{})
	urlStr := fmt.Sprintf("%s/payroll/EmploymentContracts?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	employmentContracts := []EmploymentContract{}

	for urlStr != "" {
		ac := []EmploymentContract{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetEmploymentContractsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		employmentContracts = append(employmentContracts, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &employmentContracts, nil
}

func (eo *ExactOnline) GetEmploymentContracts() (*[]EmploymentContract, error) {
	acc, err := eo.GetEmploymentContractsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
