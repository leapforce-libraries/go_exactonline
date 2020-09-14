package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// SalesInvoiceLine stores SalesInvoiceLine from exactonline
//
type SalesInvoiceLine struct {
	ID                      types.GUID  `json:"ID"`
	AmountDC                float64     `json:"AmountDC"`
	AmountFC                float64     `json:"AmountFC"`
	CostCenter              string      `json:"CostCenter"`
	CostCenterDescription   string      `json:"CostCenterDescription"`
	CostUnit                string      `json:"CostUnit"`
	CostUnitDescription     string      `json:"CostUnitDescription"`
	CustomerItemCode        string      `json:"CustomerItemCode"`
	DeliveryDate            *types.Date `json:"DeliveryDate"`
	Description             string      `json:"Description"`
	Discount                float64     `json:"Discount"`
	Division                int32       `json:"Division"`
	Employee                types.GUID  `json:"Employee"`
	EmployeeFullName        string      `json:"EmployeeFullName"`
	EndTime                 *types.Date `json:"EndTime"`
	ExtraDutyAmountFC       float64     `json:"ExtraDutyAmountFC"`
	ExtraDutyPercentage     float64     `json:"ExtraDutyPercentage"`
	GLAccount               types.GUID  `json:"GLAccount"`
	GLAccountDescription    string      `json:"GLAccountDescription"`
	InvoiceID               types.GUID  `json:"InvoiceID"`
	Item                    types.GUID  `json:"Item"`
	ItemCode                string      `json:"ItemCode"`
	ItemDescription         string      `json:"ItemDescription"`
	LineNumber              int32       `json:"LineNumber"`
	NetPrice                float64     `json:"NetPrice"`
	Notes                   string      `json:"Notes"`
	Pricelist               types.GUID  `json:"Pricelist"`
	PricelistDescription    string      `json:"PricelistDescription"`
	Project                 types.GUID  `json:"Project"`
	ProjectDescription      string      `json:"ProjectDescription"`
	ProjectWBS              types.GUID  `json:"ProjectWBS"`
	ProjectWBSDescription   string      `json:"ProjectWBSDescription"`
	Quantity                float64     `json:"Quantity"`
	SalesOrder              types.GUID  `json:"SalesOrder"`
	SalesOrderLine          types.GUID  `json:"SalesOrderLine"`
	SalesOrderLineNumber    int32       `json:"SalesOrderLineNumber"`
	SalesOrderNumber        int32       `json:"SalesOrderNumber"`
	StartTime               *types.Date `json:"StartTime"`
	Subscription            types.GUID  `json:"Subscription"`
	SubscriptionDescription string      `json:"SubscriptionDescription"`
	TaxSchedule             types.GUID  `json:"TaxSchedule"`
	TaxScheduleCode         string      `json:"TaxScheduleCode"`
	TaxScheduleDescription  string      `json:"TaxScheduleDescription"`
	UnitCode                string      `json:"UnitCode"`
	UnitDescription         string      `json:"UnitDescription"`
	UnitPrice               float64     `json:"UnitPrice"`
	VATAmountDC             float64     `json:"VATAmountDC"`
	VATAmountFC             float64     `json:"VATAmountFC"`
	VATCode                 string      `json:"VATCode"`
	VATCodeDescription      string      `json:"VATCodeDescription"`
	VATPercentage           float64     `json:"VATPercentage"`
}

func (eo *ExactOnline) GetSalesInvoiceLinesInternal(filter string) (*[]SalesInvoiceLine, error) {
	selectFields := utilities.GetTaggedFieldNames("json", SalesInvoiceLine{})
	urlStr := fmt.Sprintf("%s/bulk/SalesInvoice/SalesInvoiceLines?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	salesInvoiceLines := []SalesInvoiceLine{}

	for urlStr != "" {
		ac := []SalesInvoiceLine{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetSalesInvoiceLinesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		salesInvoiceLines = append(salesInvoiceLines, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &salesInvoiceLines, nil
}

func (eo *ExactOnline) GetSalesInvoiceLines() (*[]SalesInvoiceLine, error) {
	acc, err := eo.GetSalesInvoiceLinesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
