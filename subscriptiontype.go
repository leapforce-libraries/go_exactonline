package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Subscription stores Subscription from exactonline
//
type SubscriptionType struct {
	ID               types.GUID  `json:"ID"`
	Code             string      `json:"Code"`
	Created          *types.Date `json:"Created"`
	Creator          types.GUID  `json:"Creator"`
	CreatorFullName  string      `json:"CreatorFullName"`
	Description      string      `json:"Description"`
	Division         int32       `json:"Division"`
	Modified         *types.Date `json:"Modified"`
	Modifier         types.GUID  `json:"Modifier"`
	ModifierFullName string      `json:"ModifierFullName"`
}

func (eo *ExactOnline) GetSubscriptionTypes() *errortools.Error {
	selectFields := utilities.GetTaggedFieldNames("json", SubscriptionType{})
	urlStr := fmt.Sprintf("%s/subscription/SubscriptionTypes?$select=%s", eo.baseURL(), selectFields)
	//fmt.Println(urlStr)

	eo.SubscriptionTypes = []SubscriptionType{}

	for urlStr != "" {
		st := []SubscriptionType{}

		str, e := eo.Get(urlStr, &st)
		if e != nil {
			return e
		}

		eo.SubscriptionTypes = append(eo.SubscriptionTypes, st...)

		urlStr = str
	}

	return nil
}
