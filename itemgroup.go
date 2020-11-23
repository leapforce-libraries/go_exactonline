package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// ItemGroup stores ItemGroup from exactonline
//
type ItemGroup struct {
	ID                             types.GUID  `json:"ID"`
	Code                           string      `json:"Code"`
	Created                        *types.Date `json:"Created"`
	Creator                        types.GUID  `json:"Creator"`
	CreatorFullName                string      `json:"CreatorFullName"`
	Description                    string      `json:"Description"`
	Division                       int32       `json:"Division"`
	GLCosts                        types.GUID  `json:"GLCosts"`
	GLCostsCode                    string      `json:"GLCostsCode"`
	GLCostsDescription             string      `json:"GLCostsDescription"`
	GLPurchaseAccount              types.GUID  `json:"GLPurchaseAccount"`
	GLPurchaseAccountCode          string      `json:"GLPurchaseAccountCode"`
	GLPurchaseAccountDescription   string      `json:"GLPurchaseAccountDescription"`
	GLPurchasePriceDifference      types.GUID  `json:"GLPurchasePriceDifference"`
	GLPurchasePriceDifferenceCode  string      `json:"GLPurchasePriceDifferenceCode"`
	GLPurchasePriceDifferenceDescr string      `json:"GLPurchasePriceDifferenceDescr"`
	GLRevenue                      types.GUID  `json:"GLRevenue"`
	GLRevenueCode                  string      `json:"GLRevenueCode"`
	GLRevenueDescription           string      `json:"GLRevenueDescription"`
	GLStock                        types.GUID  `json:"GLStock"`
	GLStockCode                    string      `json:"GLStockCode"`
	GLStockDescription             string      `json:"GLStockDescription"`
	GLStockVariance                types.GUID  `json:"GLStockVariance"`
	GLStockVarianceCode            string      `json:"GLStockVarianceCode"`
	GLStockVarianceDescription     string      `json:"GLStockVarianceDescription"`
	IsDefault                      byte        `json:"IsDefault"`
	Modified                       *types.Date `json:"Modified"`
	Modifier                       types.GUID  `json:"Modifier"`
	ModifierFullName               string      `json:"ModifierFullName"`
	Notes                          string      `json:"Notes"`
}

func (eo *ExactOnline) GetItemGroupsInternal(filter string) (*[]ItemGroup, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", ItemGroup{})
	urlStr := fmt.Sprintf("%s/logistics/ItemGroups?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	itemGroups := []ItemGroup{}

	for urlStr != "" {
		its := []ItemGroup{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetItemGroupsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		itemGroups = append(itemGroups, its...)

		urlStr = str
		//urlStr = ""
	}

	return &itemGroups, nil
}

func (eo *ExactOnline) GetItemGroups() (*[]ItemGroup, *errortools.Error) {
	acc, err := eo.GetItemGroupsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}

/*
func (eo *ExactOnline) GetItemGroups() error {
	selectFields := GetJsonTaggedFieldNames(ItemGroup{})
	urlStr := fmt.Sprintf("%s%s/logistics/ItemGroups?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	//fmt.Println(urlStr)

	for urlStr != "" {
		it := []ItemGroup{}

		str, err := eo.Get(urlStr, &it)
		if err != nil {
			fmt.Println("ERROR in GetItemGroups:", err)
			fmt.Println("url:", urlStr)
			return err
		}

		eo.ItemGroups = append(eo.ItemGroups, it...)

		urlStr = str
	}

	return nil
}*/
