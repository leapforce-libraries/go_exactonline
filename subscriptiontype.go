package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Subscription stores Subscription from exactonline
//
type SubscriptionType struct {
	ID   types.GUID `json:"ID"`
	Code string     `json:"Code"`
}

func (eo *ExactOnline) GetSubscriptionTypes() error {
	selectFields := utilities.GetTaggedFieldNames("json", SubscriptionType{})
	urlStr := fmt.Sprintf("%s/subscription/SubscriptionTypes?$select=%s", eo.baseURL(), selectFields)
	//fmt.Println(urlStr)

	eo.SubscriptionTypes = []SubscriptionType{}

	for urlStr != "" {
		st := []SubscriptionType{}

		str, err := eo.Get(urlStr, &st)
		if err != nil {
			return err
		}

		eo.SubscriptionTypes = append(eo.SubscriptionTypes, st...)

		urlStr = str
	}

	return nil
}
