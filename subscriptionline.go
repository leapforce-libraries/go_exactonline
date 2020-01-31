package exactonline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	types "github.com/Leapforce-nl/go_types"
)

// SubscriptionLine stores SubscriptionLine from exactonline
//
type SubscriptionLine struct {
	ID       types.GUID  `json:"ID"`
	EntryID  types.GUID  `json:"EntryID"`
	Item     types.GUID  `json:"Item"`
	FromDate *types.Date `json:"FromDate"`
	ToDate   *types.Date `json:"ToDate"`
	UnitCode string      `json:"UnitCode"`
}

type SubscriptionLineInsert struct {
	ID       types.GUID  `json:"-"`
	EntryID  types.GUID  `json:"EntryID"`
	Item     types.GUID  `json:"Item"`
	FromDate *types.Date `json:"FromDate"`
	ToDate   *types.Date `json:"ToDate"`
	UnitCode string      `json:"UnitCode"`
}

type SubscriptionLineInsertWithSubscription struct {
	ID       types.GUID  `json:"-"`
	Item     types.GUID  `json:"Item"`
	FromDate *types.Date `json:"FromDate"`
	ToDate   *types.Date `json:"ToDate"`
	UnitCode string      `json:"UnitCode"`
}

type SubscriptionLineUpdate struct {
	ID       types.GUID  `json:"-"`
	Item     types.GUID  `json:"Item"`
	FromDate *types.Date `json:"FromDate"`
	ToDate   *types.Date `json:"ToDate"`
	UnitCode string      `json:"UnitCode"`
}

func (eo *ExactOnline) GetSubscriptionLinesInternal(filter string) (*[]SubscriptionLine, error) {
	selectFields := GetJsonTaggedFieldNames(SubscriptionLine{})
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	subscriptionlines := []SubscriptionLine{}

	for urlStr != "" {
		sl := []SubscriptionLine{}

		str, err := eo.Get(urlStr, &sl)
		if err != nil {
			fmt.Println("ERROR in GetSubscriptionLinesInternal:", err)
			fmt.Println("url:", urlStr)
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
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), s.ID.String())

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
		fmt.Println("ERROR in UpdateSubscriptionLine:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", slu)
		return err
	}

	err = eo.PutBytes(urlStr, b)
	if err != nil {
		return err
	}

	fmt.Println("\nUPDATED SubscriptionLine")
	fmt.Println("url:", urlStr)
	fmt.Println("data:", slu)

	return nil
}

// InsertSubscriptionLine inserts Subscription in ExactOnline
//
func (eo *ExactOnline) InsertSubscriptionLine(sl *SubscriptionLineInsert) error {
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines", eo.ApiUrl, strconv.Itoa(eo.Division))

	b, err := json.Marshal(sl)
	if err != nil {
		return err
	}

	type HasID struct {
		ID types.GUID `json:"ID"`
	}

	he := HasID{}

	//fmt.Println(sl)

	err = eo.PostBytes(urlStr, b, &he)
	if err != nil {
		fmt.Println("ERROR in InsertSubscriptionLine:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", sl)
		return err
	}

	fmt.Println("\nINSERTED SubscriptionLine", he.ID)
	fmt.Println("url:", urlStr)
	fmt.Println("data:", sl)
	sl.ID = he.ID

	return nil
}

// DeleteSubscription deletes Subscription in ExactOnline
//
func (eo *ExactOnline) DeleteSubscriptionLine(sl *SubscriptionLine) error {
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionLines(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), sl.ID.String())

	fmt.Println("\nDELETED SubscriptionLine", urlStr, sl.ID)

	err := eo.Delete(urlStr)
	if err != nil {
		fmt.Println("ERROR in DeleteSubscriptionLine:", err)
		fmt.Println("url:", urlStr)
		return err
	}

	fmt.Println("\nDELETED SubscriptionLine")
	fmt.Println("url:", urlStr)

	return nil
}

func (s *SubscriptionLine) FromDateTime() time.Time {

	if s.FromDate != nil {
		if (s.FromDate.Time != time.Time{}) {
			return s.FromDate.Time
		}
	}

	t, _ := time.Parse("2006-01-02", "1800-01-01")
	return t
}

func (s *SubscriptionLine) ToDateTime() time.Time {

	if s.ToDate != nil {
		if (s.ToDate.Time != time.Time{}) {
			return s.ToDate.Time
		}
	}

	t, _ := time.Parse("2006-01-02", "2099-12-31")
	return t
}
