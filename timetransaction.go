package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// TimeTransaction stores TimeTransaction from exactonline
//
type TimeTransaction struct {
	Account                 types.GUID  `json:"Account"`
	AccountName             string      `json:"AccountName"`
	Activity                types.GUID  `json:"Activity"`
	ActivityDescription     string      `json:"ActivityDescription"`
	Amount                  float64     `json:"Amount"`
	AmountFC                float64     `json:"AmountFC"`
	Attachment              types.GUID  `json:"Attachment"`
	Created                 *types.Date `json:"Created,omitempty"`
	Creator                 types.GUID  `json:"Creator"`
	CreatorFullName         string      `json:"CreatorFullName"`
	Currency                string      `json:"Currency"`
	Date                    *types.Date `json:"Date,omitempty"`
	Division                int64       `json:"Division"`
	DivisionDescription     string      `json:"DivisionDescription"`
	Employee                types.GUID  `json:"Employee"`
	EndTime                 *types.Date `json:"EndTime,omitempty"`
	EntryNumber             int64       `json:"EntryNumber"`
	ErrorText               string      `json:"ErrorText"`
	HourStatus              int64       `json:"HourStatus"`
	Item                    types.GUID  `json:"Item"`
	ItemDescription         string      `json:"ItemDescription"`
	ItemDivisable           bool        `json:"ItemDivisable"`
	Modified                *types.Date `json:"Modified,omitempty"`
	Modifier                types.GUID  `json:"Modifier"`
	ModifierFullName        string      `json:"ModifierFullName"`
	Notes                   string      `json:"Notes"`
	Price                   float64     `json:"Price"`
	PriceFC                 float64     `json:"PriceFC"`
	Project                 types.GUID  `json:"Project"`
	ProjectAccount          types.GUID  `json:"ProjectAccount"`
	ProjectAccountCode      string      `json:"ProjectAccountCode"`
	ProjectAccountName      string      `json:"ProjectAccountName"`
	ProjectCode             string      `json:"ProjectCode"`
	ProjectDescription      string      `json:"ProjectDescription"`
	Quantity                float64     `json:"Quantity"`
	StartTime               *types.Date `json:"StartTime,omitempty"`
	Subscription            types.GUID  `json:"Subscription"`
	SubscriptionAccount     types.GUID  `json:"	SubscriptionAccount"`
	SubscriptionAccountCode string      `json:"SubscriptionAccountCode"`
	SubscriptionAccountName string      `json:"SubscriptionAccountName"`
	SubscriptionDescription string      `json:"SubscriptionDescription"`
	SubscriptionNumber      int64       `json:"SubscriptionNumber"`
	Type                    int64       `json:"Type"`
}

func (eo *ExactOnline) GetTimeTransactionsInternal(filter string) (*[]TimeTransaction, error) {
	selectFields := utilities.GetTaggedFieldNames("json", TimeTransaction{})
	urlStr := fmt.Sprintf("%s/project/TimeTransactions?$select=%s", eo.baseURL(), selectFields)
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
		//urlStr = ""
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
