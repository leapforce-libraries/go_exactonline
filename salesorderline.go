package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// SalesOrderLine stores SalesOrderLine from exactonline
//
type SalesOrderLine struct {
	ID                      types.GUID  `json:"ID"`
	AmountDC                float64     `json:"AmountDC"`
	AmountFC                float64     `json:"AmountFC"`
	CostCenter              string      `json:"CostCenter"`
	CostCenterDescription   string      `json:"CostCenterDescription"`
	CostPriceFC             float64     `json:"CostPriceFC"`
	CostUnit                string      `json:"CostUnit"`
	CostUnitDescription     string      `json:"CostUnitDescription"`
	CustomerItemCode        string      `json:"CustomerItemCode"`
	DeliveryDate            *types.Date `json:"DeliveryDate"`
	Description             string      `json:"Description"`
	Discount                float64     `json:"Discount"`
	Division                int32       `json:"Division"`
	Item                    types.GUID  `json:"Item"`
	ItemCode                string      `json:"ItemCode"`
	ItemDescription         string      `json:"ItemDescription"`
	ItemVersion             types.GUID  `json:"ItemVersion"`
	ItemVersionDescription  string      `json:"ItemVersionDescription"`
	LineNumber              int32       `json:"LineNumber"`
	NetPrice                float64     `json:"NetPrice"`
	Notes                   string      `json:"Notes"`
	OrderID                 types.GUID  `json:"OrderID"`
	OrderNumber             int32       `json:"OrderNumber"`
	Pricelist               types.GUID  `json:"Pricelist"`
	PricelistDescription    string      `json:"PricelistDescription"`
	Project                 types.GUID  `json:"Project"`
	ProjectDescription      string      `json:"ProjectDescription"`
	PurchaseOrder           types.GUID  `json:"PurchaseOrder"`
	PurchaseOrderLine       types.GUID  `json:"PurchaseOrderLine"`
	PurchaseOrderLineNumber int32       `json:"PurchaseOrderLineNumber"`
	PurchaseOrderNumber     int32       `json:"PurchaseOrderNumber"`
	Quantity                float64     `json:"Quantity"`
	ShopOrder               types.GUID  `json:"ShopOrder"`
	UnitCode                string      `json:"UnitCode"`
	UnitDescription         string      `json:"UnitDescription"`
	UnitPrice               float64     `json:"UnitPrice"`
	UseDropShipment         byte        `json:"UseDropShipment"`
	VATAmount               float64     `json:"VATAmount"`
	VATCode                 string      `json:"VATCode"`
	VATCodeDescription      string      `json:"VATCodeDescription"`
	VATPercentage           float64     `json:"VATPercentage"`
}

func (eo *ExactOnline) GetSalesOrderLinesInternal(filter string) (*[]SalesOrderLine, error) {
	selectFields := GetJsonTaggedFieldNames(SalesOrderLine{})
	urlStr := fmt.Sprintf("%s%s/bulk/SalesOrder/SalesOrderLines?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	salesOrderLines := []SalesOrderLine{}

	for urlStr != "" {
		ac := []SalesOrderLine{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetSalesOrderLinesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		salesOrderLines = append(salesOrderLines, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &salesOrderLines, nil
}

func (eo *ExactOnline) GetSalesOrderLines() (*[]SalesOrderLine, error) {
	acc, err := eo.GetSalesOrderLinesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
