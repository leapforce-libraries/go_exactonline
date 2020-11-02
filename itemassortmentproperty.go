package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// ItemAssortmentProperty stores ItemAssortmentProperty from exactonline
//
type ItemAssortmentProperty struct {
	ID                 types.GUID `json:"ID "`
	Code               string     `json:"Code"`
	Description        string     `json:"Description"`
	Division           int32      `json:"Division"`
	ItemAssortmentCode int32      `json:"ItemAssortmentCode"`
}

func (eo *ExactOnline) GetItemAssortmentPropertiesInternal(filter string) (*[]ItemAssortmentProperty, error) {
	selectFields := utilities.GetTaggedFieldNames("json", ItemAssortmentProperty{})
	urlStr := fmt.Sprintf("%s/logistics/ItemAssortmentProperty?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	itemAssortmentProperties := []ItemAssortmentProperty{}

	for urlStr != "" {
		its := []ItemAssortmentProperty{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetItemAssortmentPropertiesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		itemAssortmentProperties = append(itemAssortmentProperties, its...)

		urlStr = str
		//urlStr = ""
	}

	return &itemAssortmentProperties, nil
}

func (eo *ExactOnline) GetItemAssortmentProperties() (*[]ItemAssortmentProperty, error) {
	acc, err := eo.GetItemAssortmentPropertiesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
