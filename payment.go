package exactonline

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
	utilities "github.com/Leapforce-nl/go_utilities"
)

// Payment stores Payment from exactonline
//
type Payment struct {
	ID                           types.GUID  `json:"ID"`
	Account                      types.GUID  `json:"Account"`
	AccountBankAccountID         types.GUID  `json:"AccountBankAccountID"`
	AccountBankAccountNumber     string      `json:"AccountBankAccountNumber"`
	AccountCode                  string      `json:"AccountCode"`
	AccountContact               types.GUID  `json:"AccountContact"`
	AccountContactName           string      `json:"AccountContactName"`
	AccountName                  string      `json:"AccountName"`
	AmountDC                     float64     `json:"AmountDC"`
	AmountDiscountDC             float64     `json:"AmountDiscountDC"`
	AmountDiscountFC             float64     `json:"AmountDiscountFC"`
	AmountFC                     float64     `json:"AmountFC"`
	BankAccountID                types.GUID  `json:"BankAccountID"`
	BankAccountNumber            string      `json:"BankAccountNumber"`
	CashflowTransactionBatchCode string      `json:"CashflowTransactionBatchCode"`
	Created                      *types.Date `json:"Created"`
	Creator                      types.GUID  `json:"Creator"`
	CreatorFullName              string      `json:"CreatorFullName"`
	Currency                     string      `json:"Currency"`
	Description                  string      `json:"Description"`
	DiscountDueDate              *types.Date `json:"DiscountDueDate"`
	Division                     int32       `json:"Division"`
	Document                     types.GUID  `json:"Document"`
	DocumentNumber               int32       `json:"DocumentNumber"`
	DocumentSubject              string      `json:"DocumentSubject"`
	DueDate                      *types.Date `json:"DueDate"`
	EndDate                      *types.Date `json:"EndDate"`
	EndPeriod                    int16       `json:"EndPeriod"`
	EndYear                      int16       `json:"EndYear"`
	EntryDate                    *types.Date `json:"EntryDate"`
	EntryID                      types.GUID  `json:"EntryID"`
	EntryNumber                  int32       `json:"EntryNumber"`
	GLAccount                    types.GUID  `json:"GLAccount"`
	GLAccountCode                string      `json:"GLAccountCode"`
	GLAccountDescription         string      `json:"GLAccountDescription"`
	InvoiceDate                  *types.Date `json:"InvoiceDate"`
	InvoiceNumber                int32       `json:"InvoiceNumber"`
	IsBatchBooking               byte        `json:"IsBatchBooking"`
	Journal                      string      `json:"Journal"`
	JournalDescription           string      `json:"JournalDescription"`
	Modified                     *types.Date `json:"Modified"`
	Modifier                     types.GUID  `json:"Modifier"`
	ModifierFullName             string      `json:"ModifierFullName"`
	PaymentBatchNumber           int32       `json:"PaymentBatchNumber"`
	PaymentCondition             string      `json:"PaymentCondition"`
	PaymentConditionDescription  string      `json:"PaymentConditionDescription"`
	PaymentDays                  int32       `json:"PaymentDays"`
	PaymentDaysDiscount          int32       `json:"PaymentDaysDiscount"`
	PaymentDiscountPercentage    float64     `json:"PaymentDiscountPercentage"`
	PaymentMethod                string      `json:"PaymentMethod"`
	PaymentReference             string      `json:"PaymentReference"`
	PaymentSelected              *types.Date `json:"PaymentSelected"`
	PaymentSelector              types.GUID  `json:"PaymentSelector"`
	PaymentSelectorFullName      string      `json:"PaymentSelectorFullName"`
	RateFC                       float64     `json:"RateFC"`
	Source                       int32       `json:"Source"`
	Status                       int16       `json:"Status"`
	TransactionAmountDC          float64     `json:"TransactionAmountDC"`
	TransactionAmountFC          float64     `json:"TransactionAmountFC"`
	TransactionDueDate           *types.Date `json:"TransactionDueDate"`
	TransactionEntryID           types.GUID  `json:"TransactionEntryID"`
	TransactionID                types.GUID  `json:"TransactionID"`
	TransactionIsReversal        bool        `json:"TransactionIsReversal"`
	TransactionReportingPeriod   int16       `json:"TransactionReportingPeriod"`
	TransactionReportingYear     int16       `json:"TransactionReportingYear"`
	TransactionStatus            int16       `json:"TransactionStatus"`
	TransactionType              int32       `json:"TransactionType"`
	YourRef                      string      `json:"YourRef"`
}

func (eo *ExactOnline) GetPaymentsInternal(filter string) (*[]Payment, error) {
	selectFields := utilities.GetTaggedFieldNames("json", Payment{})
	urlStr := fmt.Sprintf("%s/cashflow/Payments?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	payments := []Payment{}

	for urlStr != "" {
		its := []Payment{}

		//fmt.Println(urlStr)
		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetPaymentsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		payments = append(payments, its...)

		urlStr = str
		//urlStr = ""
	}

	return &payments, nil
}

func (eo *ExactOnline) GetPayments() (*[]Payment, error) {
	acc, err := eo.GetPaymentsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
