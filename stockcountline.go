package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// StockCountLine stores StockCountLine from exactonline
//
type StockCountLine struct {
	ID types.GUID `json:"ID"`
	//BatchNumbers       `json:"BatchNumbers"`
	CostPrice          float64     `json:"CostPrice"`
	Created            *types.Date `json:"Created"`
	Creator            types.GUID  `json:"Creator"`
	CreatorFullName    string      `json:"CreatorFullName"`
	Division           int32       `json:"Division"`
	Item               types.GUID  `json:"Item"`
	ItemCode           string      `json:"ItemCode"`
	ItemCostPrice      float64     `json:"ItemCostPrice"`
	ItemDescription    string      `json:"ItemDescription"`
	ItemDivisable      bool        `json:"ItemDivisable"`
	LineNumber         int32       `json:"LineNumber"`
	Modified           *types.Date `json:"Modified"`
	Modifier           types.GUID  `json:"Modifier"`
	ModifierFullName   string      `json:"ModifierFullName"`
	QuantityDifference float64     `json:"QuantityDifference"`
	QuantityInStock    float64     `json:"QuantityInStock"`
	QuantityNew        float64     `json:"QuantityNew"`
	//SerialNumbers              `json:"SerialNumbers"`
	StockCountID               types.GUID `json:"StockCountID"`
	StockKeepingUnit           string     `json:"StockKeepingUnit"`
	StorageLocation            types.GUID `json:"StorageLocation"`
	StorageLocationCode        string     `json:"StorageLocationCode"`
	StorageLocationDescription string     `json:"StorageLocationDescription"`
}

func (eo *ExactOnline) GetStockCountLinesInternal(filter string) (*[]StockCountLine, error) {
	selectFields := utilities.GetTaggedFieldNames("json", StockCountLine{})
	urlStr := fmt.Sprintf("%s/inventory/StockCountLines?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	stockCountLines := []StockCountLine{}

	for urlStr != "" {
		ac := []StockCountLine{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetStockCountLinesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		stockCountLines = append(stockCountLines, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &stockCountLines, nil
}

func (eo *ExactOnline) GetStockCountLines() (*[]StockCountLine, error) {
	acc, err := eo.GetStockCountLinesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
