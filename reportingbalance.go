package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// ReportingBalance stores ReportingBalance from exactonline
//
type ReportingBalance struct {
	ID                    string     `json:"ID"`
	Amount                float64    `json:"Amount"`
	AmountCredit          float64    `json:"AmountCredit"`
	AmountDebit           float64    `json:"AmountDebit"`
	BalanceType           string     `json:"BalanceType"`
	CostCenterCode        string     `json:"CostCenterCode"`
	CostCenterDescription string     `json:"CostCenterDescription"`
	CostUnitCode          string     `json:"CostUnitCode"`
	CostUnitDescription   string     `json:"CostUnitDescription"`
	Count                 int32      `json:"Count"`
	Division              int32      `json:"Division"`
	GLAccount             types.GUID `json:"GLAccount"`
	GLAccountCode         string     `json:"GLAccountCode"`
	GLAccountDescription  string     `json:"GLAccountDescription"`
	ReportingPeriod       int32      `json:"ReportingPeriod"`
	ReportingYear         int32      `json:"ReportingYear"`
	Status                int32      `json:"Status"`
	Type                  int32      `json:"Type"`
}

func (eo *ExactOnline) GetReportingBalancesInternal(filter string) (*[]ReportingBalance, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", ReportingBalance{})
	urlStr := fmt.Sprintf("%s/financial/ReportingBalance?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	reportingBalances := []ReportingBalance{}

	for urlStr != "" {
		its := []ReportingBalance{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetReportingBalancesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		reportingBalances = append(reportingBalances, its...)

		urlStr = str
		//urlStr = ""
	}

	return &reportingBalances, nil
}

func (eo *ExactOnline) GetReportingBalances() (*[]ReportingBalance, *errortools.Error) {
	acc, err := eo.GetReportingBalancesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
