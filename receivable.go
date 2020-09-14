package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// Receivable stores Receivable from exactonline
//
type Receivable struct {
	ID                            types.GUID  `json:"ID"`
	Account                       types.GUID  `json:"Account"`
	AccountBankAccountID          types.GUID  `json:"AccountBankAccountID"`
	AccountBankAccountNumber      string      `json:"AccountBankAccountNumber"`
	AccountCode                   string      `json:"AccountCode"`
	AccountContact                types.GUID  `json:"AccountContact"`
	AccountContactName            string      `json:"AccountContactName"`
	AccountCountry                string      `json:"AccountCountry"`
	AccountName                   string      `json:"AccountName"`
	AmountDC                      float64     `json:"AmountDC"`
	AmountDiscountDC              float64     `json:"AmountDiscountDC"`
	AmountDiscountFC              float64     `json:"AmountDiscountFC"`
	AmountFC                      float64     `json:"AmountFC"`
	BankAccountID                 types.GUID  `json:"BankAccountID"`
	BankAccountNumber             string      `json:"BankAccountNumber"`
	CashflowTransactionBatchCode  string      `json:"CashflowTransactionBatchCode"`
	Created                       *types.Date `json:"Created"`
	Creator                       types.GUID  `json:"Creator"`
	CreatorFullName               string      `json:"CreatorFullName"`
	Currency                      string      `json:"Currency"`
	Description                   string      `json:"Description"`
	DirectDebitMandate            types.GUID  `json:"DirectDebitMandate"`
	DirectDebitMandateDescription string      `json:"DirectDebitMandateDescription"`
	DirectDebitMandatePaymentType int16       `json:"DirectDebitMandatePaymentType"`
	DirectDebitMandateReference   string      `json:"DirectDebitMandateReference"`
	DirectDebitMandateType        int16       `json:"DirectDebitMandateType"`
	DiscountDueDate               *types.Date `json:"DiscountDueDate"`
	Division                      int32       `json:"Division"`
	Document                      types.GUID  `json:"Document"`
	DocumentNumber                int32       `json:"DocumentNumber"`
	DocumentSubject               string      `json:"DocumentSubject"`
	DueDate                       *types.Date `json:"DueDate"`
	EndDate                       *types.Date `json:"EndDate"`
	EndPeriod                     int16       `json:"EndPeriod"`
	EndToEndID                    string      `json:"EndToEndID"`
	EndYear                       int16       `json:"EndYear"`
	EntryDate                     *types.Date `json:"EntryDate"`
	EntryID                       types.GUID  `json:"EntryID"`
	EntryNumber                   int32       `json:"EntryNumber"`
	GLAccount                     types.GUID  `json:"GLAccount"`
	GLAccountCode                 string      `json:"GLAccountCode"`
	GLAccountDescription          string      `json:"GLAccountDescription"`
	InvoiceDate                   *types.Date `json:"InvoiceDate"`
	InvoiceNumber                 int32       `json:"InvoiceNumber"`
	IsBatchBooking                byte        `json:"IsBatchBooking"`
	IsFullyPaid                   bool        `json:"IsFullyPaid"`
	Journal                       string      `json:"Journal"`
	JournalDescription            string      `json:"JournalDescription"`
	LastPaymentDate               *types.Date `json:"LastPaymentDate"`
	Modified                      *types.Date `json:"Modified"`
	Modifier                      types.GUID  `json:"Modifier"`
	ModifierFullName              string      `json:"ModifierFullName"`
	PaymentCondition              string      `json:"PaymentCondition"`
	PaymentConditionDescription   string      `json:"PaymentConditionDescription"`
	PaymentDays                   int32       `json:"PaymentDays"`
	PaymentDaysDiscount           int32       `json:"PaymentDaysDiscount"`
	PaymentDiscountPercentage     float64     `json:"PaymentDiscountPercentage"`
	PaymentInformationID          string      `json:"PaymentInformationID"`
	PaymentMethod                 string      `json:"PaymentMethod"`
	PaymentReference              string      `json:"PaymentReference"`
	RateFC                        float64     `json:"RateFC"`
	ReceivableBatchNumber         int32       `json:"ReceivableBatchNumber"`
	ReceivableSelected            *types.Date `json:"ReceivableSelected"`
	ReceivableSelector            types.GUID  `json:"ReceivableSelector"`
	ReceivableSelectorFullName    string      `json:"ReceivableSelectorFullName"`
	Source                        int32       `json:"Source"`
	Status                        int16       `json:"Status"`
	TransactionAmountDC           float64     `json:"TransactionAmountDC"`
	TransactionAmountFC           float64     `json:"TransactionAmountFC"`
	TransactionDueDate            *types.Date `json:"TransactionDueDate"`
	TransactionEntryID            types.GUID  `json:"TransactionEntryID"`
	TransactionID                 types.GUID  `json:"TransactionID"`
	TransactionIsReversal         bool        `json:"TransactionIsReversal"`
	TransactionReportingPeriod    int16       `json:"TransactionReportingPeriod"`
	TransactionReportingYear      int16       `json:"TransactionReportingYear"`
	TransactionStatus             int16       `json:"TransactionStatus"`
	TransactionType               int32       `json:"TransactionType"`
	YourRef                       string      `json:"YourRef"`
}

func (eo *ExactOnline) GetReceivablesInternal(filter string) (*[]Receivable, error) {
	selectFields := utilities.GetTaggedFieldNames("json", Receivable{})
	urlStr := fmt.Sprintf("%s/cashflow/Receivables?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	receivables := []Receivable{}

	for urlStr != "" {
		ps := []Receivable{}

		str, err := eo.Get(urlStr, &ps)
		if err != nil {
			fmt.Println("ERROR in GetReceivablesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		receivables = append(receivables, ps...)

		urlStr = str
		//urlStr = ""
	}

	return &receivables, nil
}

func (eo *ExactOnline) GetReceivables() (*[]Receivable, error) {
	acc, err := eo.GetReceivablesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
