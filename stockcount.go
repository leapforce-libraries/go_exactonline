package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// StockCount stores StockCount from exactonline
//
type StockCount struct {
	StockCountID                 types.GUID  `json:"StockCountID"`
	Created                      *types.Date `json:"Created"`
	Creator                      types.GUID  `json:"Creator"`
	CreatorFullName              string      `json:"CreatorFullName"`
	Description                  string      `json:"Description"`
	Division                     int32       `json:"Division"`
	EntryNumber                  int32       `json:"EntryNumber"`
	Modified                     *types.Date `json:"Modified"`
	Modifier                     types.GUID  `json:"Modifier"`
	ModifierFullName             string      `json:"ModifierFullName"`
	OffsetGLInventory            types.GUID  `json:"OffsetGLInventory"`
	OffsetGLInventoryCode        string      `json:"OffsetGLInventoryCode"`
	OffsetGLInventoryDescription string      `json:"OffsetGLInventoryDescription"`
	Source                       int16       `json:"Source"`
	Status                       int16       `json:"Status"`
	StockCountDate               *types.Date `json:"StockCountDate"`
	//StockCountLines              `json:"StockCountLines"`
	StockCountNumber     int32      `json:"StockCountNumber"`
	Warehouse            types.GUID `json:"Warehouse"`
	WarehouseCode        string     `json:"WarehouseCode"`
	WarehouseDescription string     `json:"WarehouseDescription"`
}

func (eo *ExactOnline) GetStockCountsInternal(filter string) (*[]StockCount, error) {
	selectFields := utilities.GetTaggedFieldNames("json", StockCount{})
	urlStr := fmt.Sprintf("%s/inventory/StockCounts?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	stockCounts := []StockCount{}

	for urlStr != "" {
		ac := []StockCount{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetStockCountsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		stockCounts = append(stockCounts, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &stockCounts, nil
}

func (eo *ExactOnline) GetStockCounts() (*[]StockCount, error) {
	acc, err := eo.GetStockCountsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
