package exactonline

import (
	"fmt"

	utilities "github.com/leapforce-libraries/go_utilities"
)

// OutstandingInvoicesOverview stores OutstandingInvoicesOverview from exactonline
//
type OutstandingInvoicesOverview struct {
	CurrencyCode                       string  `json:"CurrencyCode"`
	OutstandingPayableInvoiceAmount    float64 `json:"OutstandingPayableInvoiceAmount"`
	OutstandingPayableInvoiceCount     float64 `json:"OutstandingPayableInvoiceCount"`
	OutstandingReceivableInvoiceAmount float64 `json:"OutstandingReceivableInvoiceAmount"`
	OutstandingReceivableInvoiceCount  float64 `json:"OutstandingReceivableInvoiceCount"`
	OverduePayableInvoiceAmount        float64 `json:"OverduePayableInvoiceAmount"`
	OverduePayableInvoiceCount         float64 `json:"OverduePayableInvoiceCount"`
	OverdueReceivableInvoiceAmount     float64 `json:"OverdueReceivableInvoiceAmount"`
	OverdueReceivableInvoiceCount      float64 `json:"OverdueReceivableInvoiceCount"`
}

func (eo *ExactOnline) GetOutstandingInvoicesOverviewsInternal(filter string) (*[]OutstandingInvoicesOverview, error) {
	selectFields := utilities.GetTaggedFieldNames("json", OutstandingInvoicesOverview{})
	urlStr := fmt.Sprintf("%s/read/financial/OutstandingInvoicesOverview?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	outstandingInvoicesOverviews := []OutstandingInvoicesOverview{}

	for urlStr != "" {
		ac := []OutstandingInvoicesOverview{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetOutstandingInvoicesOverviewsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		outstandingInvoicesOverviews = append(outstandingInvoicesOverviews, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &outstandingInvoicesOverviews, nil
}

func (eo *ExactOnline) GetOutstandingInvoicesOverviews() (*[]OutstandingInvoicesOverview, error) {
	acc, err := eo.GetOutstandingInvoicesOverviewsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
