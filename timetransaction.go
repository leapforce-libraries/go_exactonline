package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// TimeTransaction stores TimeTransaction from exactonline
//
type TimeTransaction struct {
	Account             types.GUID  `json:"Account"`
	AccountName         string      `json:"AccountName"`
	Activity            types.GUID  `json:"Activity"`
	ActivityDescription string      `json:"ActivityDescription"`
	Amount              float64     `json:"Amount"`
	AmountFC            float64     `json:"AmountFC"`
	Attachment          types.GUID  `json:"Attachment"`
	Created             *types.Date `json:"Created,omitempty"`
	Creator             types.GUID  `json:"Creator"`
	CreatorFullName     string      `json:"CreatorFullName"`
	Currency            string      `json:"Currency"`
	Date                *types.Date `json:"Date,omitempty"`
	Division            int64       `json:"Division"`
	DivisionDescription string      `json:"DivisionDescription"`
}

func (eo *ExactOnline) GetTimeTransactionsInternal(filter string) (*[]TimeTransaction, error) {
	selectFields := GetJsonTaggedFieldNames(TimeTransaction{})
	urlStr := fmt.Sprintf("%s%s/project/TimeTransactions?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	timeTransactions := []TimeTransaction{}

	for urlStr != "" {
		ac := []TimeTransaction{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetTimeTransactionsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		timeTransactions = append(timeTransactions, ac...)

		urlStr = str

		urlStr = ""
	}

	return &timeTransactions, nil
}

func (eo *ExactOnline) GetTimeTransactions() (*[]TimeTransaction, error) {
	acc, err := eo.GetTimeTransactionsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
