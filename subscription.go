package exactonline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	types "types"
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
	StartDate                      types.Date `json:"StartDate,omitempty"`
	EndDate                        types.Date `json:"EndDate,omitempty"`
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

// ToBq convert Subscription to SubscriptionBq
//
func (s *Subscription) IsValid(timestamp time.Time) bool {
	if s.StartDate != *new(types.Date) {
		if s.StartDate.After(timestamp) {
			return false
		}
	}

	if s.EndDate != *new(types.Date) {
		if s.EndDate.Before(timestamp) {
			return false
		}
	}

	return true
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
		//fmt.Println("len(sub.SubscriptionLines)", len(account.Subscriptions[i].SubscriptionLines))
		//fmt.Println("sd/ed", account.Subscriptions[i].StartDate, account.Subscriptions[i].EndDate)
	}

	//fmt.Println("GetSubscriptionsByAccount:", len(account.Subscriptions))
	return nil
}

type SubscriptionUpdate struct {
	SubscriptionType types.GUID `json:"SubscriptionType"`
	OrderedBy        types.GUID `json:"OrderedBy"`
	InvoiceTo        types.GUID `json:"InvoiceTo"`
	StartDate        types.Date `json:"StartDate"`
	EndDate          types.Date `json:"EndDate"`
	Description      string     `json:"Description"`
}

type SubscriptionInsert struct {
	SubscriptionType  types.GUID         `json:"SubscriptionType"`
	OrderedBy         types.GUID         `json:"OrderedBy"`
	InvoiceTo         types.GUID         `json:"InvoiceTo"`
	StartDate         types.Date         `json:"StartDate"`
	EndDate           types.Date         `json:"EndDate"`
	Description       string             `json:"Description"`
	SubscriptionLines []SubscriptionLine `json:"SubscriptionLines"`
}

// UpdateSubscription updates Subscription in ExactOnline
//
func (eo *ExactOnline) UpdateSubscription(s *Subscription) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), s.EntryID.String())

	/*sd := new(types.Date)
	if !s.StartDate.IsZero() {
		sd = &s.StartDate
	}
	ed := new(types.Date)
	if !s.EndDate.IsZero() {
		ed = &s.EndDate
	}*/
	su := SubscriptionUpdate{
		s.SubscriptionType,
		s.OrderedBy,
		s.InvoiceTo,
		s.StartDate,
		s.EndDate,
		s.Description + "_updated",
	}

	b, err := json.Marshal(su)
	if err != nil {
		return err
	}

	fmt.Println("\nUPDATED Subscription", urlStr, su)

	err = eo.PutBytes(urlStr, b)
	if err != nil {
		return err
	}

	return nil
}

// InsertSubscription inserts Subscription in ExactOnline
//
func (eo *ExactOnline) InsertSubscription(s *Subscription) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision))

	/*sd := new(types.Date)
	if !s.StartDate.IsZero() {
		sd = &s.StartDate
	}
	ed := new(types.Date)
	if !s.EndDate.IsZero() {
		ed = &s.EndDate
	}*/
	si := SubscriptionInsert{
		s.SubscriptionType,
		s.OrderedBy,
		s.InvoiceTo,
		s.StartDate,
		s.EndDate,
		s.Description + "_inserted",
		s.SubscriptionLines,
	}

	b, err := json.Marshal(si)
	if err != nil {
		return err
	}

	type HasEntryID struct {
		EntryID types.GUID `json:"EntryID"`
	}

	he := HasEntryID{}

	fmt.Println("\nINSERTED Subscription", urlStr, si)

	err = eo.PostBytes(urlStr, b, &he)
	if err != nil {
		return err
	}

	fmt.Println("\nNEW Subscription", he.EntryID)
	s.EntryID = he.EntryID

	return nil
}
