package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// ItemAssortment stores ItemAssortment from exactonline
//
type ItemAssortment struct {
	ID          types.GUID `json:"ID"`
	Code        int32      `json:"Code"`
	Description string     `json:"Description"`
	Division    int32      `json:"Division"`
	//Properties  []ItemAssortmentProperty `json:"Properties"`
}

func (eo *ExactOnline) GetItemAssortmentsInternal(filter string) (*[]ItemAssortment, error) {
	selectFields := utilities.GetTaggedFieldNames("json", ItemAssortment{})
	urlStr := fmt.Sprintf("%s/logistics/ItemAssortment?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	itemAssortments := []ItemAssortment{}

	for urlStr != "" {
		its := []ItemAssortment{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetItemAssortmentsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		itemAssortments = append(itemAssortments, its...)

		urlStr = str
		//urlStr = ""
	}

	return &itemAssortments, nil
}

func (eo *ExactOnline) GetItemAssortments() (*[]ItemAssortment, error) {
	acc, err := eo.GetItemAssortmentsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
