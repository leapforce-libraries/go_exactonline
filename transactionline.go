package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// TransactionLine stores TransactionLine from exactonline
//
type TransactionLine struct {
	Account               types.GUID  `json:"Account"`
	AccountCode           string      `json:"AccountCode"`
	AccountName           string      `json:"AccountName"`
	AmountDC              float64     `json:"AmountDC"`
	AmountFC              float64     `json:"AmountFC"`
	AmountVATBaseFC       float64     `json:"AmountVATBaseFC"`
	AmountVATFC           float64     `json:"AmountVATFC"`
	Asset                 types.GUID  `json:"Asset"`
	AssetCode             string      `json:"AssetCode"`
	AssetDescription      string      `json:"AssetDescription"`
	CostCenter            string      `json:"CostCenter"`
	CostCenterDescription string      `json:"CostCenterDescription"`
	CostUnit              string      `json:"CostUnit"`
	CostUnitDescription   string      `json:"CostUnitDescription"`
	Created               *types.Date `json:"Created,omitempty"`
	Creator               types.GUID  `json:"Creator"`
	CreatorFullName       string      `json:"CreatorFullName"`
	Currency              string      `json:"Currency"`
	Date                  *types.Date `json:"Date,omitempty"`
	Description           string      `json:"Description"`
	Division              int64       `json:"Division"`
	Document              types.GUID  `json:"Document"`
}

func (eo *ExactOnline) GetTransactionLinesInternal(filter string) (*[]TransactionLine, error) {
	selectFields := GetJsonTaggedFieldNames(TransactionLine{})
	urlStr := fmt.Sprintf("%s%s/financialtransaction/TransactionLines?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	transactionLines := []TransactionLine{}

	for urlStr != "" {
		ac := []TransactionLine{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetTransactionLinesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		transactionLines = append(transactionLines, ac...)

		urlStr = str

		urlStr = ""
	}

	return &transactionLines, nil
}

func (eo *ExactOnline) GetTransactionLines() (*[]TransactionLine, error) {
	acc, err := eo.GetTransactionLinesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
