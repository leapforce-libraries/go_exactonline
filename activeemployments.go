package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
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
	Employee                 types.GUID  `json:"Employee"`
	EmployeeFullName         string      `json:"EmployeeFullName"`
	EmployeeHID              int64       `json:"EmployeeHID"`
	EmploymentOrganization   types.GUID  `json:"EmploymentOrganization"`
	EndDate                  *types.Date `json:"EndDate,omitempty"`
	HID                      int64       `json:"HID"`
	HourlyWage               float64     `json:"HourlyWage"`
	InternalRate             float64     `json:"InternalRate"`
	Jobtitle                 types.GUID  `json:"Jobtitle"`
	JobtitleDescription      string      `json:"JobtitleDescription"`
	Modified                 *types.Date `json:"Modified,omitempty"`
	Modifier                 types.GUID  `json:"Modifier"`
	ModifierFullName         string      `json:"ModifierFullName"`
	ReasonEnd                int64       `json:"ReasonEnd"`
	ReasonEndDescription     string      `json:"ReasonEndDescription"`
	ReasonEndFlex            int64       `json:"ReasonEndFlex"`
	ReasonEndFlexDescription string      `json:"ReasonEndFlexDescription"`
	Salary                   types.GUID  `json:"Salary"`
	Schedule                 types.GUID  `json:"Schedule"`
	ScheduleAverageHours     float64     `json:"ScheduleAverageHours"`
	ScheduleCode             string      `json:"ScheduleCode"`
	ScheduleDays             float64     `json:"ScheduleDays"`
	ScheduleDescription      string      `json:"ScheduleDescription"`
	ScheduleHours            float64     `json:"ScheduleHours"`
	StartDate                *types.Date `json:"StartDate,omitempty"`
	StartDateOrganization    *types.Date `json:"StartDateOrganization,omitempty"`
}

func (eo *ExactOnline) GetActiveEmploymentsInternal(filter string) (*[]ActiveEmployment, error) {
	selectFields := utilities.GetTaggedFieldNames("json", ActiveEmployment{})
	urlStr := fmt.Sprintf("%s/payroll/ActiveEmployments?$select=%s", eo.baseURL(), selectFields)
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
		//urlStr = ""
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
