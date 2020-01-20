package exactonline

import (
	"fmt"
	"strconv"
	"time"

	"types"
)

// Subscription stores Subscription from exactonline
//
type Subscription struct {
	EntryID                        types.GUID `json:"EntryID"`
	Description                    string     `json:"Description"`
	Division                       int        `json:"Division"`
	InvoiceTo                      types.GUID `json:"InvoiceTo"`
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
	EntryID                        string
	Description                    string
	Division                       int
	InvoiceTo                      string
	InvoiceToContactPersonFullName string
	InvoiceToName                  string
	OrderedBy                      string
	SubscriptionTypeCode           string
	StartDate                      time.Time
	EndDate                        time.Time
	//SubscriptionLines              []SubscriptionLine `json:"SubscriptionLines"`
}

// ToBq convert Subscription to SubscriptionBq
//
func (s *Subscription) ToBq() *SubscriptionBq {
	return &SubscriptionBq{
		s.EntryID.String(),
		s.Description,
		s.Division,
		s.InvoiceTo.String(),
		s.InvoiceToContactPersonFullName,
		s.InvoiceToName,
		s.OrderedBy.String(),
		s.SubscriptionTypeCode,
		s.StartDate.Time,
		s.EndDate.Time,
		//s.SubscriptionLines,
	}
}

func (eo *ExactOnline) getSubscriptions() error {
	selectFields := GetJsonTaggedFieldNames(Subscription{})
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
	//fmt.Println(urlStr)

	eo.Subscriptions = []Subscription{}

	for urlStr != "" {
		su := []Subscription{}

		str, err := eo.get(urlStr, &su)
		if err != nil {
			return err
		}

		eo.Subscriptions = append(eo.Subscriptions, su...)
		//fmt.Println(len(eo.Accounts))

		urlStr = str
		//urlStr = "" //temp
	}

	return nil
}

func (eo *ExactOnline) UpdateSubscription(s *Subscription) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), s.EntryID.String())

	data := make(map[string]string)
	data["SubscriptionTypeCode"] = s.SubscriptionTypeCode

	err := eo.put(urlStr, data)
	if err != nil {
		return err
	}

	fmt.Println("updated SubscriptionTypeCode", urlStr, s.SubscriptionTypeCode)

	//time.Sleep(1 * time.Second)

	return nil
}
