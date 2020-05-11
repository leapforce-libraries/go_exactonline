package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// Item stores Item from exactonline
//
type Item struct {
	ID          types.GUID `json:"ID"`
	Code        string     `json:"Code"`
	Description string     `json:"Description"`
}

func (eo *ExactOnline) GetItemsInternal(filter string) (*[]Item, error) {
	selectFields := GetJsonTaggedFieldNames(Item{})
	urlStr := fmt.Sprintf("%s%s/logistics/Items?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	items := []Item{}

	for urlStr != "" {
		its := []Item{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetItemsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		items = append(items, its...)

		urlStr = str
		//urlStr = ""
	}

	return &items, nil
}

func (eo *ExactOnline) GetItems() (*[]Item, error) {
	acc, err := eo.GetItemsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}

/*
func (eo *ExactOnline) GetItems() error {
	selectFields := GetJsonTaggedFieldNames(Item{})
	urlStr := fmt.Sprintf("%s%s/logistics/Items?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	//fmt.Println(urlStr)

	for urlStr != "" {
		it := []Item{}

		str, err := eo.Get(urlStr, &it)
		if err != nil {
			fmt.Println("ERROR in GetItems:", err)
			fmt.Println("url:", urlStr)
			return err
		}

		eo.Items = append(eo.Items, it...)

		urlStr = str
	}

	return nil
}*/
