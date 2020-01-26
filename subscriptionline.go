package exactonline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	types "github.com/leapforce-nl/go_types"
)

// SubscriptionLine stores SubscriptionLine from exactonline
//
type SubscriptionLine struct {
	ID       types.GUID `json:"ID"`
	EntryID  types.GUID `json:"EntryID"`
	Item     types.GUID `json:"Item"`
	FromDate types.Date `json:"FromDate"`
	ToDate   types.Date `json:"ToDate"`
	UnitCode string     `json:"UnitCode"`
}

type SubscriptionLineInsert struct {
	ID       types.GUID `json:"-"`
	EntryID  types.GUID `json:"EntryID"`
	Item     types.GUID `json:"Item"`
	FromDate types.Date `json:"FromDate"`
	ToDate   types.Date `json:"ToDate"`
	UnitCode string     `json:"UnitCode"`
}

type SubscriptionLineUpdate struct {
	ID       types.GUID `json:"ID"`
	Item     types.GUID `json:"Item"`
	FromDate types.Date `json:"FromDate"`
	ToDate   types.Date `json:"ToDate"`
	UnitCode string     `json:"UnitCode"`
}

func (eo *ExactOnline) GetSubscriptionLinesInternal(filter string) (*[]SubscriptionLine, error) {
	selectFields := GetJsonTaggedFieldNames(SubscriptionLine{})
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	subscriptionlines := []SubscriptionLine{}

	for urlStr != "" {
		sl := []SubscriptionLine{}

		str, err := eo.Get(urlStr, &sl)
		if err != nil {
			return nil, err
		}

		subscriptionlines = append(subscriptionlines, sl...)

		urlStr = str
	}

	return &subscriptionlines, nil
}

func (eo *ExactOnline) GetSubscriptionLines() error {
	sub, err := eo.GetSubscriptionLinesInternal("")
	if err != nil {
		return err
	}
	eo.SubscriptionLines = *sub

	return nil
}

// GetSubscriptionLinesBySubscription return all SubscriptionLines for a single Subscription
//
func (eo ExactOnline) GetSubscriptionLinesBySubscription(subscription *Subscription) error {
	filter := fmt.Sprintf("EntryID eq guid'%s'", subscription.EntryID.String())

	sub, err := eo.GetSubscriptionLinesInternal(filter)
	if err != nil {
		return err
	}
	subscription.SubscriptionLines = *sub

	//fmt.Println("GetSubscriptionLinesBySubscription:", len(subscription.SubscriptionLines))
	return nil
}

// UpdateSubscription updates Subscription in ExactOnline
//
func (eo *ExactOnline) UpdateSubscriptionLine(s *SubscriptionLine) error {
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), s.EntryID.String())

	/*sd := new(types.Date)
	if !s.StartDate.IsZero() {
		sd = &s.StartDate
	}
	ed := new(types.Date)
	if !s.EndDate.IsZero() {
		ed = &s.EndDate
	}*/
	slu := SubscriptionLineUpdate{
		s.ID,
		s.Item,
		s.FromDate,
		s.ToDate,
		s.UnitCode,
	}

	b, err := json.Marshal(slu)
	if err != nil {
		return err
	}

	fmt.Println("\nUPDATED SubscriptionLine", urlStr, slu)

	err = eo.PutBytes(urlStr, b)
	if err != nil {
		return err
	}

	return nil
}

// InsertSubscriptionLine inserts Subscription in ExactOnline
//
func (eo *ExactOnline) InsertSubscriptionLine(sl *SubscriptionLineInsert) error {
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision))

	b, err := json.Marshal(sl)
	if err != nil {
		return err
	}

	type HasID struct {
		ID types.GUID `json:"ID"`
	}

	he := HasID{}

	//fmt.Println(sl)

	fmt.Println("\nINSERTED SubscriptionLine", urlStr, sl)

	err = eo.PostBytes(urlStr, b, &he)
	if err != nil {
		return err
	}

	fmt.Println("\nNEW SubscriptionLine", he.ID)
	sl.ID = he.ID

	return nil
}

// DeleteSubscription deletes Subscription in ExactOnline
//
func (eo *ExactOnline) DeleteSubscriptionLine(sl *SubscriptionLine) error {
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), sl.EntryID.String())

	fmt.Println("\nDELETED SubscriptionLine", urlStr, sl.EntryID)

	err := eo.Delete(urlStr)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionLine) FromDateTime() time.Time {
	if (s.FromDate.Time == time.Time{}) {
		t, _ := time.Parse("2006-01-02", "1800-01-01")
		return t
	}
	return s.FromDate.Time
}

func (s *SubscriptionLine) ToDateTime() time.Time {
	if (s.ToDate.Time == time.Time{}) {
		t, _ := time.Parse("2006-01-02", "2099-12-31")
		return t
	}
	return s.ToDate.Time
}
