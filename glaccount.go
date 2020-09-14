package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// GLAccount stores GLAccount from exactonline
//
type GLAccount struct {
	ID                             types.GUID  `json:"ID"`
	AssimilatedVATBox              int16       `json:"AssimilatedVATBox"`
	BalanceSide                    string      `json:"BalanceSide"`
	BalanceType                    string      `json:"BalanceType"`
	BelcotaxType                   int32       `json:"BelcotaxType"`
	Code                           string      `json:"Code"`
	Compress                       bool        `json:"Compress"`
	Costcenter                     string      `json:"Costcenter"`
	CostcenterDescription          string      `json:"CostcenterDescription"`
	Costunit                       string      `json:"Costunit"`
	CostunitDescription            string      `json:"CostunitDescription"`
	Created                        *types.Date `json:"Created"`
	Creator                        types.GUID  `json:"Creator"`
	CreatorFullName                string      `json:"CreatorFullName"`
	Description                    string      `json:"Description"`
	Division                       int32       `json:"Division"`
	ExcludeVATListing              byte        `json:"ExcludeVATListing"`
	ExpenseNonDeductiblePercentage float64     `json:"ExpenseNonDeductiblePercentage"`
	IsBlocked                      bool        `json:"IsBlocked"`
	Matching                       bool        `json:"Matching"`
	Modified                       *types.Date `json:"Modified"`
	Modifier                       types.GUID  `json:"Modifier"`
	ModifierFullName               string      `json:"ModifierFullName"`
	PrivateGLAccount               types.GUID  `json:"PrivateGLAccount"`
	PrivatePercentage              float64     `json:"PrivatePercentage"`
	ReportingCode                  string      `json:"ReportingCode"`
	RevalueCurrency                bool        `json:"RevalueCurrency"`
	SearchCode                     string      `json:"SearchCode"`
	Type                           int32       `json:"Type"`
	TypeDescription                string      `json:"TypeDescription"`
	UseCostcenter                  byte        `json:"UseCostcenter"`
	UseCostunit                    byte        `json:"UseCostunit"`
	VATCode                        string      `json:"VATCode"`
	VATDescription                 string      `json:"VATDescription"`
	VATGLAccountType               string      `json:"VATGLAccountType"`
	VATNonDeductibleGLAccount      types.GUID  `json:"VATNonDeductibleGLAccount"`
	VATNonDeductiblePercentage     float64     `json:"VATNonDeductiblePercentage"`
	VATSystem                      string      `json:"VATSystem"`
	YearEndCostGLAccount           types.GUID  `json:"YearEndCostGLAccount"`
	YearEndReflectionGLAccount     types.GUID  `json:"YearEndReflectionGLAccount"`
}

func (eo *ExactOnline) GetGLAccountsInternal(filter string) (*[]GLAccount, error) {
	selectFields := utilities.GetTaggedFieldNames("json", GLAccount{})
	urlStr := fmt.Sprintf("%s/bulk/financial/GLAccounts?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	glAccounts := []GLAccount{}

	for urlStr != "" {
		ac := []GLAccount{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetGLAccountsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		glAccounts = append(glAccounts, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &glAccounts, nil
}

func (eo *ExactOnline) GetGLAccounts() (*[]GLAccount, error) {
	acc, err := eo.GetGLAccountsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
