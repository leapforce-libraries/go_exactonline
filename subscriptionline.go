package exactonline

import (
	"encoding/json"
	"fmt"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// SubscriptionLine stores SubscriptionLine from exactonline
//
type SubscriptionLine struct {
	ID       types.GUID `json:"ID"`
	EntryID  types.GUID `json:"EntryID"`
	Item     types.GUID `json:"Item"`
	ItemCode string
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

func (s *SubscriptionLine) FromDateString() string {
	if s.FromDate == nil {
		return "-"
	}
	return s.FromDate.Time.Format("2006-01-02")
}

func (s *SubscriptionLine) ToDateString() string {
	if s.ToDate == nil {
		return "-"
	}
	return s.ToDate.Time.Format("2006-01-02")
}

/*
var oldSubscriptionLine *SubscriptionLine
var newSubscriptionLine *SubscriptionLine

// SaveValues saves current values in local copy of Contact
//
func (s *SubscriptionLine) SaveValues(inserted bool) {
	oldSubscriptionLine = nil
	if !inserted {
		oldSubscriptionLine = new(SubscriptionLine)
		oldSubscriptionLine.ItemCode = s.ItemCode
		oldSubscriptionLine.FromDate = s.FromDate
		oldSubscriptionLine.ToDate = s.ToDate
		oldSubscriptionLine.UnitCode = s.UnitCode
	}
}

func (s *SubscriptionLine) Values(deleted bool) (string, string) {
	oldValues := ""
	newValues := ""

	newSubscriptionLine = nil
	if !deleted {
		newSubscriptionLine = new(SubscriptionLine)
		newSubscriptionLine.ItemCode = s.ItemCode
		newSubscriptionLine.FromDate = s.FromDate
		newSubscriptionLine.ToDate = s.ToDate
		newSubscriptionLine.UnitCode = s.UnitCode
	}

	if oldSubscriptionLine == nil {
		newValues += ",ItemCode:" + newSubscriptionLine.ItemCode
	} else if newSubscriptionLine == nil {
		oldValues += ",ItemCode:" + oldSubscriptionLine.ItemCode
	} else if oldSubscriptionLine.ItemCode != newSubscriptionLine.ItemCode {
		oldValues += ",ItemCode:" + oldSubscriptionLine.ItemCode
		newValues += ",ItemCode:" + newSubscriptionLine.ItemCode
	}

	if oldSubscriptionLine == nil {
		newValues += ",FromDate:" + newSubscriptionLine.FromDateString()
	} else if newSubscriptionLine == nil {
		oldValues += ",FromDate:" + oldSubscriptionLine.FromDateString()
	} else if oldSubscriptionLine.FromDateString() != newSubscriptionLine.FromDateString() {
		oldValues += ",FromDate:" + oldSubscriptionLine.FromDateString()
		newValues += ",FromDate:" + newSubscriptionLine.FromDateString()
	}

	if oldSubscriptionLine == nil {
		newValues += ",ToDate:" + newSubscriptionLine.ToDateString()
	} else if newSubscriptionLine == nil {
		oldValues += ",ToDate:" + oldSubscriptionLine.ToDateString()
	} else if oldSubscriptionLine.ToDateString() != newSubscriptionLine.ToDateString() {
		oldValues += ",ToDate:" + oldSubscriptionLine.ToDateString()
		newValues += ",ToDate:" + newSubscriptionLine.ToDateString()
	}

	if oldSubscriptionLine == nil {
		newValues += ",UnitCode:" + newSubscriptionLine.UnitCode
	} else if newSubscriptionLine == nil {
		oldValues += ",UnitCode:" + oldSubscriptionLine.UnitCode
	} else if oldSubscriptionLine.UnitCode != newSubscriptionLine.UnitCode {
		oldValues += ",UnitCode:" + oldSubscriptionLine.UnitCode
		newValues += ",UnitCode:" + newSubscriptionLine.UnitCode
	}

	oldValues = strings.TrimLeft(oldValues, ",")
	newValues = strings.TrimLeft(newValues, ",")

	return oldValues, newValues
}*/

func (eo *ExactOnline) GetSubscriptionLinesInternal(filter string) (*[]SubscriptionLine, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", SubscriptionLine{})
	urlStr := fmt.Sprintf("%s/subscription/SubscriptionLines?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	subscriptionlines := []SubscriptionLine{}

	for urlStr != "" {
		sl := []SubscriptionLine{}

		str, e := eo.Get(urlStr, &sl)
		if e != nil {
			return nil, e
		}

		for ii := range sl {
			for _, item := range eo.Items {
				if sl[ii].Item == item.ID {
					sl[ii].ItemCode = item.Code
					//fmt.Println("sl[ii].ItemCode", sl[ii].ItemCode)
					break
				}
			}
		}

		subscriptionlines = append(subscriptionlines, sl...)

		urlStr = str
	}

	return &subscriptionlines, nil
}

func (eo *ExactOnline) GetSubscriptionLines() *errortools.Error {
	sub, e := eo.GetSubscriptionLinesInternal("")
	if e != nil {
		return e
	}
	eo.SubscriptionLines = *sub

	return nil
}

// GetSubscriptionLinesBySubscription return all SubscriptionLines for a single Subscription
//
func (eo ExactOnline) GetSubscriptionLinesBySubscription(subscription *Subscription) *errortools.Error {
	filter := fmt.Sprintf("EntryID eq guid'%s'", subscription.EntryID.String())

	sub, e := eo.GetSubscriptionLinesInternal(filter)
	if e != nil {
		return e
	}
	subscription.SubscriptionLines = *sub

	//fmt.Println("GetSubscriptionLinesBySubscription:", len(subscription.SubscriptionLines))
	return nil
}

// UpdateSubscription updates Subscription in ExactOnline
//
func (eo *ExactOnline) UpdateSubscriptionLine(s *SubscriptionLine) *errortools.Error {
	urlStr := fmt.Sprintf("%s/subscription/SubscriptionLines(guid'%s')", eo.baseURL(), s.ID.String())

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
		return errortools.ErrorMessage(err)
	}

	e := eo.PutBytes(urlStr, b)
	if e != nil {
		return e
	}

	return nil
}

// InsertSubscriptionLine inserts Subscription in ExactOnline
//
func (eo *ExactOnline) InsertSubscriptionLine(sl *SubscriptionLine) *errortools.Error {
	if sl == nil {
		return nil
	}

	urlStr := fmt.Sprintf("%s/subscription/SubscriptionLines", eo.baseURL())

	sli := SubscriptionLineInsert{}
	sli.EntryID = sl.EntryID
	sli.Item = sl.Item
	sli.FromDate = sl.FromDate
	sli.ToDate = sl.ToDate
	sli.UnitCode = sl.UnitCode

	b, err := json.Marshal(sli)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	type HasID struct {
		ID types.GUID `json:"ID"`
	}

	he := HasID{}

	//fmt.Println(sl)

	err_ := eo.PostBytes(urlStr, b, &he)
	if err_ != nil {
		return err_
	}

	sl.ID = he.ID

	return nil
}

// DeleteSubscription deletes Subscription in ExactOnline
//
func (eo *ExactOnline) DeleteSubscriptionLine(sl *SubscriptionLine) *errortools.Error {
	if sl == nil {
		return nil
	}

	urlStr := fmt.Sprintf("%s/subscription/SubscriptionLines(guid'%s')", eo.baseURL(), sl.ID.String())

	fmt.Println("\nDELETED SubscriptionLine", urlStr, sl.ID)

	e := eo.Delete(urlStr)
	if e != nil {
		return e
	}

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
