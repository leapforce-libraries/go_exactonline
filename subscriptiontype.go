package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/leapforce-nl/go_types"
)

// Subscription stores Subscription from exactonline
//
type SubscriptionType struct {
	ID   types.GUID `json:"ID"`
	Code string     `json:"Code"`
}

func (eo *ExactOnline) GetSubscriptionTypes() error {
	selectFields := GetJsonTaggedFieldNames(SubscriptionType{})
	urlStr := fmt.Sprintf("%s%s/subscription/SubscriptionTypes?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
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
