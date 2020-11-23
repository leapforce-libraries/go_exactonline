package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// SalesOrder stores SalesOrder from exactonline
//
type SalesOrder struct {
	OrderID                        types.GUID  `json:"OrderID"`
	AmountDC                       float64     `json:"AmountDC"`
	AmountDiscount                 float64     `json:"AmountDiscount"`
	AmountDiscountExclVat          float64     `json:"AmountDiscountExclVat"`
	AmountFC                       float64     `json:"AmountFC"`
	AmountFCExclVat                float64     `json:"AmountFCExclVat"`
	ApprovalStatus                 int16       `json:"ApprovalStatus"`
	ApprovalStatusDescription      string      `json:"ApprovalStatusDescription"`
	Approved                       *types.Date `json:"Approved"`
	Approver                       types.GUID  `json:"Approver"`
	ApproverFullName               string      `json:"ApproverFullName"`
	Created                        *types.Date `json:"Created"`
	Creator                        types.GUID  `json:"Creator"`
	CreatorFullName                string      `json:"CreatorFullName"`
	Currency                       string      `json:"Currency"`
	DeliverTo                      types.GUID  `json:"DeliverTo"`
	DeliverToContactPerson         types.GUID  `json:"DeliverToContactPerson"`
	DeliverToContactPersonFullName string      `json:"DeliverToContactPersonFullName"`
	DeliverToName                  string      `json:"DeliverToName"`
	DeliveryAddress                types.GUID  `json:"DeliveryAddress"`
	DeliveryDate                   *types.Date `json:"DeliveryDate"`
	DeliveryStatus                 int16       `json:"DeliveryStatus"`
	DeliveryStatusDescription      string      `json:"DeliveryStatusDescription"`
	Description                    string      `json:"Description"`
	Discount                       float64     `json:"Discount"`
	Division                       int32       `json:"Division"`
	Document                       types.GUID  `json:"Document"`
	DocumentNumber                 int32       `json:"DocumentNumber"`
	DocumentSubject                string      `json:"DocumentSubject"`
	InvoiceStatus                  int16       `json:"InvoiceStatus"`
	InvoiceStatusDescription       string      `json:"InvoiceStatusDescription"`
	InvoiceTo                      types.GUID  `json:"InvoiceTo"`
	InvoiceToContactPerson         types.GUID  `json:"InvoiceToContactPerson"`
	InvoiceToContactPersonFullName string      `json:"InvoiceToContactPersonFullName"`
	InvoiceToName                  string      `json:"InvoiceToName"`
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
	//SalesOrderLines                `json:"SalesOrderLines"`
	Salesperson               types.GUID `json:"Salesperson"`
	SalespersonFullName       string     `json:"SalespersonFullName"`
	SelectionCode             types.GUID `json:"SelectionCode"`
	SelectionCodeCode         string     `json:"SelectionCodeCode"`
	SelectionCodeDescription  string     `json:"SelectionCodeDescription"`
	ShippingMethod            types.GUID `json:"ShippingMethod"`
	ShippingMethodDescription string     `json:"ShippingMethodDescription"`
	Status                    int16      `json:"Status"`
	StatusDescription         string     `json:"StatusDescription"`
	TaxSchedule               types.GUID `json:"TaxSchedule"`
	TaxScheduleCode           string     `json:"TaxScheduleCode"`
	TaxScheduleDescription    string     `json:"TaxScheduleDescription"`
	WarehouseCode             string     `json:"WarehouseCode"`
	WarehouseDescription      string     `json:"WarehouseDescription"`
	WarehouseID               types.GUID `json:"WarehouseID"`
	YourRef                   string     `json:"YourRef"`
}

func (eo *ExactOnline) GetSalesOrdersInternal(filter string) (*[]SalesOrder, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", SalesOrder{})
	urlStr := fmt.Sprintf("%s/bulk/SalesOrder/SalesOrders?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	salesOrders := []SalesOrder{}

	for urlStr != "" {
		ac := []SalesOrder{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetSalesOrdersInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		salesOrders = append(salesOrders, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &salesOrders, nil
}

func (eo *ExactOnline) GetSalesOrders() (*[]SalesOrder, *errortools.Error) {
	acc, err := eo.GetSalesOrdersInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
