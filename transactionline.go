package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// TransactionLine stores TransactionLine from exactonline
//
type TransactionLine struct {
	Account                   types.GUID  `json:"Account"`
	AccountCode               string      `json:"AccountCode"`
	AccountName               string      `json:"AccountName"`
	AmountDC                  float64     `json:"AmountDC"`
	AmountFC                  float64     `json:"AmountFC"`
	AmountVATBaseFC           float64     `json:"AmountVATBaseFC"`
	AmountVATFC               float64     `json:"AmountVATFC"`
	Asset                     types.GUID  `json:"Asset"`
	AssetCode                 string      `json:"AssetCode"`
	AssetDescription          string      `json:"AssetDescription"`
	CostCenter                string      `json:"CostCenter"`
	CostCenterDescription     string      `json:"CostCenterDescription"`
	CostUnit                  string      `json:"CostUnit"`
	CostUnitDescription       string      `json:"CostUnitDescription"`
	Created                   *types.Date `json:"Created,omitempty"`
	Creator                   types.GUID  `json:"Creator"`
	CreatorFullName           string      `json:"CreatorFullName"`
	Currency                  string      `json:"Currency"`
	Date                      *types.Date `json:"Date,omitempty"`
	Description               string      `json:"Description"`
	Division                  int64       `json:"Division"`
	Document                  types.GUID  `json:"Document"`
	DocumentNumber            int64       `json:"DocumentNumber"`
	DocumentSubject           string      `json:"DocumentSubject"`
	DueDate                   *types.Date `json:"DueDate,omitempty"`
	EntryID                   types.GUID  `json:"EntryID"`
	EntryNumber               int64       `json:"EntryNumber"`
	ExchangeRate              float64     `json:"ExchangeRate"`
	ExtraDutyAmountFC         float64     `json:"ExtraDutyAmountFC"`
	ExtraDutyPercentage       float64     `json:"ExtraDutyPercentage"`
	FinancialPeriod           int64       `json:"FinancialPeriod"`
	FinancialYear             int64       `json:"FinancialYear"`
	GLAccount                 types.GUID  `json:"GLAccount"`
	GLAccountCode             string      `json:"GLAccountCode"`
	GLAccountDescription      string      `json:"GLAccountDescription"`
	InvoiceNumber             int64       `json:"InvoiceNumber"`
	Item                      types.GUID  `json:"Item"`
	ItemCode                  string      `json:"ItemCode"`
	ItemDescription           string      `json:"ItemDescription"`
	JournalCode               string      `json:"JournalCode"`
	JournalDescription        string      `json:"JournalDescription"`
	LineNumber                int64       `json:"LineNumber"`
	LineType                  int64       `json:"LineType"`
	Modified                  *types.Date `json:"Modified,omitempty"`
	Modifier                  types.GUID  `json:"Modifier"`
	ModifierFullName          string      `json:"ModifierFullName"`
	Notes                     string      `json:"Notes"`
	OffsetID                  types.GUID  `json:"OffsetID"`
	OrderNumber               int64       `json:"OrderNumber"`
	PaymentDiscountAmount     float64     `json:"PaymentDiscountAmount"`
	PaymentReference          string      `json:"PaymentReference"`
	Project                   types.GUID  `json:"Project"`
	ProjectCode               string      `json:"ProjectCode"`
	ProjectDescription        string      `json:"ProjectDescription"`
	Quantity                  float64     `json:"Quantity"`
	SerialNumber              string      `json:"SerialNumber"`
	Status                    int64       `json:"Status"`
	Subscription              types.GUID  `json:"Subscription"`
	SubscriptionDescription   string      `json:"SubscriptionDescription"`
	TrackingNumber            string      `json:"TrackingNumber"`
	TrackingNumberDescription string      `json:"TrackingNumberDescription"`
	Type                      int64       `json:"Type"`
	VATCode                   string      `json:"VATCode"`
	VATCodeDescription        string      `json:"VATCodeDescription"`
	VATPercentage             float64     `json:"VATPercentage"`
	VATType                   string      `json:"VATType"`
	YourRef                   string      `json:"YourRef"`
}

func (eo *ExactOnline) GetTransactionLinesInternal(filter string) (*[]TransactionLine, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", TransactionLine{})
	urlStr := fmt.Sprintf("%s/bulk/financial/TransactionLines?$select=%s", eo.baseURL(), selectFields)
	//urlStr := fmt.Sprintf("%s%s/financialtransaction/TransactionLines?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
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
		//urlStr = ""
	}

	return &transactionLines, nil
}

func (eo *ExactOnline) GetTransactionLines() (*[]TransactionLine, *errortools.Error) {
	acc, err := eo.GetTransactionLinesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
