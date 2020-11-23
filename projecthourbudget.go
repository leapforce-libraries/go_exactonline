package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// ProjectHourBudget stores ProjectHourBudget from exactonline
//
type ProjectHourBudget struct {
	ID                 types.GUID  `json:"ID"`
	Budget             float64     `json:"Budget"`
	Created            *types.Date `json:"Created,omitempty"`
	Creator            types.GUID  `json:"Creator"`
	CreatorFullName    string      `json:"CreatorFullName"`
	Division           int64       `json:"Division"`
	Item               types.GUID  `json:"Item"`
	ItemCode           string      `json:"ItemCode"`
	ItemDescription    string      `json:"ItemDescription"`
	Modified           *types.Date `json:"Modified,omitempty"`
	Modifier           types.GUID  `json:"Modifier"`
	ModifierFullName   string      `json:"ModifierFullName"`
	Project            types.GUID  `json:"Project"`
	ProjectCode        string      `json:"ProjectCode"`
	ProjectDescription string      `json:"ProjectDescription"`
}

func (eo *ExactOnline) GetProjectHourBudgetsInternal(filter string) (*[]ProjectHourBudget, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", ProjectHourBudget{})
	urlStr := fmt.Sprintf("%s/project/ProjectHourBudgets?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	projectHourBudgets := []ProjectHourBudget{}

	for urlStr != "" {
		ac := []ProjectHourBudget{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetProjectHourBudgetsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		projectHourBudgets = append(projectHourBudgets, ac...)

		urlStr = str
		//urlStr = ""
	}

	/*if len(projectHourBudgets) > 0 {
		fmt.Println("#ProjectHourBudget:", len(projectHourBudgets))
	}*/

	return &projectHourBudgets, nil
}

func (eo *ExactOnline) GetProjectHourBudgets() (*[]ProjectHourBudget, *errortools.Error) {
	acc, err := eo.GetProjectHourBudgetsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (eo *ExactOnline) GetProjectHourBudgetsByProject(projectID types.GUID) (*[]ProjectHourBudget, *errortools.Error) {
	filter := fmt.Sprintf("Project eq guid'%s'", projectID.String())

	acc, err := eo.GetProjectHourBudgetsInternal(filter)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
