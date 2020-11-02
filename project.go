package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Project stores Project from exactonline
//
type Project struct {
	ID                        types.GUID `json:"ID"`
	Account                   types.GUID `json:"Account"`
	AccountCode               string     `json:"AccountCode"`
	AccountContact            types.GUID `json:"AccountContact"`
	AccountName               string     `json:"AccountName"`
	AllowAdditionalInvoicing  bool       `json:"AllowAdditionalInvoicing"`
	BlockEntry                bool       `json:"BlockEntry"`
	BlockRebilling            bool       `json:"BlockRebilling"`
	BudgetedAmount            float64    `json:"BudgetedAmount"`
	BudgetedCosts             float64    `json:"BudgetedCosts"`
	BudgetedHoursPerHourType  []ProjectHourBudget
	BudgetedRevenue           float64     `json:"BudgetedRevenue"`
	BudgetOverrunHours        byte        `json:"BudgetOverrunHours"`
	BudgetType                int64       `json:"BudgetType"`
	BudgetTypeDescription     string      `json:"BudgetTypeDescription"`
	Classification            types.GUID  `json:"Classification"`
	ClassificationDescription string      `json:"ClassificationDescription"`
	Code                      string      `json:"Code"`
	CostsAmountFC             float64     `json:"CostsAmountFC"`
	Created                   *types.Date `json:"Created,omitempty"`
	Creator                   types.GUID  `json:"Creator"`
	CreatorFullName           string      `json:"CreatorFullName"`
	CustomerPOnumber          string      `json:"CustomerPOnumber"`
	Description               string      `json:"Description"`
	Division                  int64       `json:"Division"`
	DivisionName              string      `json:"DivisionName"`
	EndDate                   *types.Date `json:"EndDate,omitempty"`
	FixedPriceItem            types.GUID  `json:"FixedPriceItem"`
	FixedPriceItemDescription string      `json:"FixedPriceItemDescription"`
	HasWBSLines               bool        `json:"HasWBSLines"`
	InternalNotes             string      `json:"InternalNotes"`
	InvoiceAsQuoted           bool        `json:"InvoiceAsQuoted"`
	Manager                   types.GUID  `json:"Manager"`
	ManagerFullname           string      `json:"ManagerFullname"`
	MarkupPercentage          float64     `json:"MarkupPercentage"`
	Modified                  *types.Date `json:"Modified,omitempty"`
	Modifier                  types.GUID  `json:"Modifier"`
	ModifierFullName          string      `json:"ModifierFullName"`
	Notes                     string      `json:"Notes"`
	PrepaidItem               types.GUID  `json:"PrepaidItem"`
	PrepaidItemDescription    string      `json:"PrepaidItemDescription"`
	PrepaidType               int64       `json:"PrepaidType"`
	PrepaidTypeDescription    string      `json:"PrepaidTypeDescription"`
	SalesTimeQuantity         float64     `json:"SalesTimeQuantity"`
	SourceQuotation           types.GUID  `json:"SourceQuotation"`
	StartDate                 *types.Date `json:"StartDate,omitempty"`
	TimeQuantityToAlert       float64     `json:"TimeQuantityToAlert"`
	Type                      int64       `json:"Type"`
	TypeDescription           string      `json:"TypeDescription"`
	UseBillingMilestones      bool        `json:"UseBillingMilestones"`
}

func (eo *ExactOnline) GetProjectsInternal(filter string) (*[]Project, error) {
	selectFields := utilities.GetTaggedFieldNames("json", Project{})
	urlStr := fmt.Sprintf("%s/project/Projects?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	projects := []Project{}

	for urlStr != "" {
		ps := []Project{}

		str, err := eo.Get(urlStr, &ps)
		if err != nil {
			fmt.Println("ERROR in GetProjectsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		for index, p := range ps {
			budgetedHours, err := eo.GetProjectHourBudgetsByProject(p.ID)
			if err != nil {
				return nil, err
			}

			ps[index].BudgetedHoursPerHourType = *budgetedHours
		}

		projects = append(projects, ps...)

		urlStr = str
		//urlStr = ""
	}

	return &projects, nil
}

func (eo *ExactOnline) GetProjects() (*[]Project, error) {
	acc, err := eo.GetProjectsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
