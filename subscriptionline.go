package exactonline

import (
	"encoding/json"
	"fmt"
	"strconv"

	types "github.com/leapforce-nl/go_types"
)

// SubscriptionLine stores SubscriptionLine from exactonline
//
type SubscriptionLine struct {
	EntryID  types.GUID `json:"EntryID"`
	Item     types.GUID `json:"Item"`
	FromDate types.Date `json:"FromDate"`
	ToDate   types.Date `json:"ToDate"`
	UnitCode string     `json:"UnitCode"`
}

type SubscriptionLineInsert struct {
	EntryID  types.GUID `json:"-"`
	Item     types.GUID `json:"Item"`
	FromDate types.Date `json:"FromDate"`
	ToDate   types.Date `json:"ToDate"`
	UnitCode string     `json:"UnitCode"`
}

type SubscriptionLineUpdate struct {
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
