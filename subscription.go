package exactonline

import (
	"time"

	"github.com/mcnijman/go-exactonline/types"
)

// Subscription stores Subscription from exactonline
//
type Subscription struct {
	Description                    string     `json:"Description"`
	Division                       int        `json:"Division"`
	InvoiceToContactPersonFullName string     `json:"InvoiceToContactPersonFullName"`
	InvoiceToName                  string     `json:"InvoiceToName"`
	OrderedBy                      types.GUID `json:"OrderedBy"`
	SubscriptionTypeCode           string     `json:"SubscriptionTypeCode"`
	StartDate                      types.Date `json:"StartDate"`
	EndDate                        types.Date `json:"EndDate"`
	//SubscriptionLines              []SubscriptionLine `json:"SubscriptionLines"`
}

// SubscriptionBq equals type Subscription except fields of type Date that are converted to type Time, to be insertable in BigQuery
//
type SubscriptionBq struct {
	Description                    string
	Division                       int
	InvoiceToContactPersonFullName string
	InvoiceToName                  string
	OrderedBy                      types.GUID
	SubscriptionTypeCode           string
	StartDate                      time.Time
	EndDate                        time.Time
	//SubscriptionLines              []SubscriptionLine `json:"SubscriptionLines"`
}

// ToBq convert Subscription to SubscriptionBq
//
func (s *Subscription) ToBq() *SubscriptionBq {
	return &SubscriptionBq{
		s.Description,
		s.Division,
		s.InvoiceToContactPersonFullName,
		s.InvoiceToName,
		s.OrderedBy,
		s.SubscriptionTypeCode,
		s.StartDate.Time,
		s.EndDate.Time,
		//s.SubscriptionLines,
	}
}
