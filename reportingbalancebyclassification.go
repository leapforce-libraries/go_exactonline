package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// ReportingBalanceByClassification stores ReportingBalanceByClassification from exactonline
//
type ReportingBalanceByClassification struct {
	ID                        string     `json:"ID"`
	Amount                    float64    `json:"Amount"`
	AmountCredit              float64    `json:"AmountCredit"`
	AmountDebit               float64    `json:"AmountDebit"`
	BalanceType               string     `json:"BalanceType"`
	ClassificationCode        string     `json:"ClassificationCode"`
	ClassificationDescription string     `json:"ClassificationDescription"`
	CostCenterCode            string     `json:"CostCenterCode"`
	CostCenterDescription     string     `json:"CostCenterDescription"`
	CostUnitCode              string     `json:"CostUnitCode"`
	CostUnitDescription       string     `json:"CostUnitDescription"`
	Count                     int32      `json:"Count"`
	Division                  int32      `json:"Division"`
	GLAccount                 types.GUID `json:"GLAccount"`
	GLAccountCode             string     `json:"GLAccountCode"`
	GLAccountDescription      string     `json:"GLAccountDescription"`
	GLScheme                  types.GUID `json:"GLScheme"`
	ReportingPeriod           int32      `json:"ReportingPeriod"`
	ReportingYear             int32      `json:"ReportingYear"`
	Status                    int32      `json:"Status"`
	Type                      int32      `json:"Type"`
}

func (eo *ExactOnline) GetReportingBalanceByClassificationsInternal(glScheme GLScheme, reportingYear int, filter string) (*[]ReportingBalanceByClassification, error) {
	selectFields := utilities.GetTaggedFieldNames("json", ReportingBalanceByClassification{})
	urlStr := fmt.Sprintf("%s/read/financial/ReportingBalanceByClassification?glScheme=guid'%s'&reportingYear=%s&$select=%s", eo.baseURL(), glScheme.ID.String(), strconv.Itoa(reportingYear), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	reportingBalanceByClassifications := []ReportingBalanceByClassification{}

	for urlStr != "" {
		ac := []ReportingBalanceByClassification{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetReportingBalanceByClassificationsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		reportingBalanceByClassifications = append(reportingBalanceByClassifications, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &reportingBalanceByClassifications, nil
}

func (eo *ExactOnline) GetReportingBalanceByClassifications(glScheme GLScheme, reportingYear int) (*[]ReportingBalanceByClassification, error) {
	acc, err := eo.GetReportingBalanceByClassificationsInternal(glScheme, reportingYear, "")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
