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
	SubscriptionType               types.GUID `json:"SubscriptionType"`
	SubscriptionTypeCode           string     `json:"SubscriptionTypeCode"`
	StartDate                      types.Date `json:"StartDate"`
	EndDate                        types.Date `json:"EndDate"`
	SubscriptionLines              []SubscriptionLine
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

func (eo *ExactOnline) GetSubscriptionsInternal(filter string) (*[]Subscription, error) {
	selectFields := GetJsonTaggedFieldNames(Subscription{})
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	subscriptions := []Subscription{}

	for urlStr != "" {
		sc := []Subscription{}

		str, err := eo.Get(urlStr, &sc)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sc...)

		urlStr = str
	}

	return &subscriptions, nil
}

func (eo *ExactOnline) GetSubscriptions() error {
	sub, err := eo.GetSubscriptionsInternal("")
	if err != nil {
		return err
	}
	eo.Subscriptions = *sub

	return nil
}

// GetSubscriptionsByAccount return all Subscriptions for a single Account
//
func (eo ExactOnline) GetSubscriptionsByAccount(account *Account) error {
	filter := fmt.Sprintf("OrderedBy eq guid'%s'", account.ID.String())

	sub, err := eo.GetSubscriptionsInternal(filter)
	if err != nil {
		return err
	}

	account.Subscriptions = *sub

	for i := range account.Subscriptions {
		err = eo.GetSubscriptionLinesBySubscription(&account.Subscriptions[i])
		if err != nil {
			return err
		}
		fmt.Println("len(sub.SubscriptionLines)", len(account.Subscriptions[i].SubscriptionLines))
	}

	fmt.Println("GetSubscriptionsByAccount:", len(account.Subscriptions))
	return nil
}

// UpdateSubscription updates Subscription in ExactOnline
//
func (eo *ExactOnline) UpdateSubscription(s *Subscription) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), s.EntryID.String())

	data := make(map[string]string)
	data["SubscriptionType"] = s.SubscriptionType.String()

	err := eo.Put(urlStr, data)
	if err != nil {
		return err
	}

	fmt.Println("updated SubscriptionType", urlStr, s.SubscriptionType)

	//time.Sleep(1 * time.Second)

	return nil
}
