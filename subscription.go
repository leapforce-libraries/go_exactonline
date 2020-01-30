package exactonline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	types "github.com/Leapforce-nl/go_types"
)

// Subscription stores Subscription from exactonline
//
type Subscription struct {
	EntryID     types.GUID `json:"EntryID"`
	Description string     `json:"Description"`
	Division    int        `json:"Division"`
	//InvoiceTo                      types.GUID  `json:"InvoiceTo"`
	//InvoiceToContactPersonFullName string      `json:"InvoiceToContactPersonFullName"`
	//InvoiceToName                  string      `json:"InvoiceToName"`
	OrderedBy            types.GUID  `json:"OrderedBy"`
	SubscriptionType     types.GUID  `json:"SubscriptionType"`
	SubscriptionTypeCode string      `json:"SubscriptionTypeCode"`
	StartDate            *types.Date `json:"StartDate,omitempty"`
	EndDate              *types.Date `json:"EndDate,omitempty"`
	CancellationDate     *types.Date `json:"CancellationDate,omitempty"`
	SubscriptionLines    []SubscriptionLine
}

type SubscriptionBigQuery struct {
	EntryID              string
	Description          string
	OrderedBy            string
	SubscriptionType     string
	SubscriptionTypeCode string
	StartDate            string
	EndDate              string
	CancellationDate     string
}

type SubscriptionUpdate struct {
	SubscriptionType types.GUID `json:"SubscriptionType"`
	OrderedBy        types.GUID `json:"OrderedBy"`
	//InvoiceTo        types.GUID  `json:"InvoiceTo"`
	StartDate        *types.Date `json:"StartDate"`
	CancellationDate *types.Date `json:"CancellationDate"`
	//EndDate          types.Date `json:"EndDate"`
	Description string `json:"Description"`
}

type SubscriptionInsert struct {
	EntryID          types.GUID `json:"-"`
	SubscriptionType types.GUID `json:"SubscriptionType"`
	OrderedBy        types.GUID `json:"OrderedBy"`
	//InvoiceTo        types.GUID  `json:"InvoiceTo"`
	StartDate        *types.Date `json:"StartDate"`
	CancellationDate *types.Date `json:"CancellationDate"`
	//EndDate           types.Date                               `json:"EndDate"`
	Description       string                                   `json:"Description"`
	SubscriptionLines []SubscriptionLineInsertWithSubscription `json:"SubscriptionLines"`
}

// IsValid returns whether or not a Subscription is valid at a certain time.Time
//
func (s *Subscription) IsValid(timestamp time.Time) bool {
	if s.StartDate != new(types.Date) {
		if s.StartDate.After(timestamp) {
			return false
		}
	}

	if s.EndDate != new(types.Date) {
		if s.EndDate.Before(timestamp) {
			return false
		}
	}

	return true
}

func (s *Subscription) ToBigQuery() *SubscriptionBigQuery {
	startDate := ""
	endDate := ""
	cancellationDate := ""

	if s.StartDate != nil {
		startDate = s.StartDate.Time.Format("2006-01-02")
	}
	if s.EndDate != nil {
		endDate = s.EndDate.Time.Format("2006-01-02")
	}
	if s.CancellationDate != nil {
		cancellationDate = s.CancellationDate.Time.Format("2006-01-02")
	}

	return &SubscriptionBigQuery{
		s.EntryID.String(),
		s.Description,
		s.OrderedBy.String(),
		s.SubscriptionType.String(),
		s.SubscriptionTypeCode,
		startDate,
		endDate,
		cancellationDate,
	}
}

func (eo *ExactOnline) CancellationDate(endDate *types.Date) *types.Date {
	cancellationDate := endDate
	if cancellationDate != nil {
		if int(cancellationDate.Month()) == 12 && cancellationDate.Day() == 31 {
			// in order to set EndDate at 31/12/year we need to set the CancellationDate
			// to 30/11/year (or earlier)......
			cd, _ := time.Parse("2006-01-02", strconv.Itoa(cancellationDate.Time.Year())+"-11-30")
			cancellationDate = &types.Date{cd}
		}
	}

	return cancellationDate
}

func (eo *ExactOnline) GetSubscriptionsInternal(filter string) (*[]Subscription, error) {
	selectFields := GetJsonTaggedFieldNames(Subscription{})
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	subscriptions := []Subscription{}

	for urlStr != "" {
		sc := []Subscription{}

		str, err := eo.Get(urlStr, &sc)
		if err != nil {
			fmt.Println("ERROR in GetSubscriptionsInternal:", err)
			fmt.Println("url:", urlStr)
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

// UpdateSubscription updates Subscription in ExactOnline
//
func (eo *ExactOnline) UpdateSubscription(s *Subscription) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), s.EntryID.String())

	su := SubscriptionUpdate{
		s.SubscriptionType,
		s.OrderedBy,
		//s.InvoiceTo,
		s.StartDate,
		eo.CancellationDate(s.EndDate),
		s.Description,
	}

	b, err := json.Marshal(su)
	if err != nil {
		fmt.Println("ERROR in UpdateSubscription:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", su)
		return err
	}
	err = eo.PutBytes(urlStr, b)
	if err != nil {
		return err
	}

	fmt.Println("\nUPDATED Subscription")
	fmt.Println("url:", urlStr)
	fmt.Println("data:", su)

	return nil
}

// InsertSubscription inserts Subscription in ExactOnline
//
func (eo *ExactOnline) InsertSubscription(s *SubscriptionInsert) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions", eo.ApiUrl, strconv.Itoa(eo.Division))

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	type HasEntryID struct {
		EntryID types.GUID `json:"EntryID"`
	}

	he := HasEntryID{}

	err = eo.PostBytes(urlStr, b, &he)
	if err != nil {
		fmt.Println("ERROR in InsertSubscription:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", s)
		return err
	}

	fmt.Println("\nINSERTED Subscription", he.EntryID)
	fmt.Println("url:", urlStr)
	fmt.Println("data:", s)
	s.EntryID = he.EntryID

	return nil
}

// DeleteSubscription deletes Subscription in ExactOnline
//
func (eo *ExactOnline) DeleteSubscription(s *Subscription) error {
	urlStr := fmt.Sprintf("%s%s/subscription/Subscriptions(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), s.EntryID.String())

	err := eo.Delete(urlStr)
	if err != nil {
		fmt.Println("ERROR in DeleteSubscription:", err)
		fmt.Println("url:", urlStr)
		return err
	}

	fmt.Println("\nDELETED Subscription", urlStr)

	return nil
}
