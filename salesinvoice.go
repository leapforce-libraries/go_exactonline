package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// SalesInvoice stores SalesInvoice from exactonline
//
type SalesInvoice struct {
	InvoiceID                      types.GUID  `json:"InvoiceID"`
	AmountDC                       float64     `json:"AmountDC"`
	AmountDiscount                 float64     `json:"AmountDiscount"`
	AmountDiscountExclVat          float64     `json:"AmountDiscountExclVat"`
	AmountFC                       float64     `json:"AmountFC"`
	AmountFCExclVat                float64     `json:"AmountFCExclVat"`
	Created                        *types.Date `json:"Created"`
	Creator                        types.GUID  `json:"Creator"`
	CreatorFullName                string      `json:"CreatorFullName"`
	Currency                       string      `json:"Currency"`
	DeliverTo                      types.GUID  `json:"DeliverTo"`
	DeliverToAddress               types.GUID  `json:"DeliverToAddress"`
	DeliverToContactPerson         types.GUID  `json:"DeliverToContactPerson"`
	DeliverToContactPersonFullName string      `json:"DeliverToContactPersonFullName"`
	DeliverToName                  string      `json:"DeliverToName"`
	Description                    string      `json:"Description"`
	Discount                       float64     `json:"Discount"`
	DiscountType                   int16       `json:"DiscountType"`
	Division                       int32       `json:"Division"`
	Document                       types.GUID  `json:"Document"`
	DocumentNumber                 int32       `json:"DocumentNumber"`
	DocumentSubject                string      `json:"DocumentSubject"`
	DueDate                        *types.Date `json:"DueDate"`
	ExtraDutyAmountFC              float64     `json:"ExtraDutyAmountFC"`
	GAccountAmountFC               float64     `json:"GAccountAmountFC"`
	InvoiceDate                    *types.Date `json:"InvoiceDate"`
	InvoiceNumber                  int32       `json:"InvoiceNumber"`
	InvoiceTo                      types.GUID  `json:"InvoiceTo"`
	InvoiceToContactPerson         types.GUID  `json:"InvoiceToContactPerson"`
	InvoiceToContactPersonFullName string      `json:"InvoiceToContactPersonFullName"`
	InvoiceToName                  string      `json:"InvoiceToName"`
	IsExtraDuty                    bool        `json:"IsExtraDuty"`
	Journal                        string      `json:"Journal"`
	JournalDescription             string      `json:"JournalDescription"`
	Modified                       *types.Date `json:"Modified"`
	Modifier                       types.GUID  `json:"Modifier"`
	ModifierFullName               string      `json:"ModifierFullName"`
	OrderDate                      *types.Date `json:"OrderDate"`
	OrderedBy                      types.GUID  `json:"OrderedBy"`
	OrderedByContactPerson         types.GUID  `json:"OrderedByContactPerson"`
	OrderedByContactPersonFullName string      `json:"OrderedByContactPersonFullName"`
	OrderedByName                  string      `json:"OrderedByName"`
	OrderNumber                    int32       `json:"OrderNumber"`
	PaymentCondition               string      `json:"PaymentCondition"`
	PaymentConditionDescription    string      `json:"PaymentConditionDescription"`
	PaymentReference               string      `json:"PaymentReference"`
	Remarks                        string      `json:"Remarks"`
	//SalesInvoiceLines                    `json:"SalesInvoiceLines"`
	Salesperson                          types.GUID `json:"Salesperson"`
	SalespersonFullName                  string     `json:"SalespersonFullName"`
	SelectionCode                        types.GUID `json:"SelectionCode"`
	SelectionCodeCode                    string     `json:"SelectionCodeCode"`
	SelectionCodeDescription             string     `json:"SelectionCodeDescription"`
	StarterSalesInvoiceStatus            int16      `json:"StarterSalesInvoiceStatus"`
	StarterSalesInvoiceStatusDescription string     `json:"StarterSalesInvoiceStatusDescription"`
	Status                               int16      `json:"Status"`
	StatusDescription                    string     `json:"StatusDescription"`
	TaxSchedule                          types.GUID `json:"TaxSchedule"`
	TaxScheduleCode                      string     `json:"TaxScheduleCode"`
	TaxScheduleDescription               string     `json:"TaxScheduleDescription"`
	Type                                 int32      `json:"Type"`
	TypeDescription                      string     `json:"TypeDescription"`
	VATAmountDC                          float64    `json:"VATAmountDC"`
	VATAmountFC                          float64    `json:"VATAmountFC"`
	Warehouse                            types.GUID `json:"Warehouse"`
	WithholdingTaxAmountFC               float64    `json:"WithholdingTaxAmountFC"`
	WithholdingTaxBaseAmount             float64    `json:"WithholdingTaxBaseAmount"`
	WithholdingTaxPercentage             float64    `json:"WithholdingTaxPercentage"`
	YourRef                              string     `json:"YourRef"`
}

func (eo *ExactOnline) GetSalesInvoicesInternal(filter string) (*[]SalesInvoice, error) {
	selectFields := utilities.GetTaggedFieldNames("json", SalesInvoice{})
	urlStr := fmt.Sprintf("%s/bulk/SalesInvoice/SalesInvoices?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	salesInvoices := []SalesInvoice{}

	for urlStr != "" {
		ac := []SalesInvoice{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetSalesInvoicesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		salesInvoices = append(salesInvoices, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &salesInvoices, nil
}

func (eo *ExactOnline) GetSalesInvoices() (*[]SalesInvoice, error) {
	acc, err := eo.GetSalesInvoicesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
