package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// Asset stores Asset from exactonline
//
type Asset struct {
	ID                            types.GUID  `json:"ID "`
	AlreadyDepreciated            byte        `json:"AlreadyDepreciated"`
	AssetFrom                     types.GUID  `json:"AssetFrom"`
	AssetFromDescription          string      `json:"AssetFromDescription"`
	AssetGroup                    types.GUID  `json:"AssetGroup"`
	AssetGroupCode                string      `json:"AssetGroupCode"`
	AssetGroupDescription         string      `json:"AssetGroupDescription"`
	CatalogueValue                float64     `json:"CatalogueValue"`
	Code                          string      `json:"Code"`
	Costcenter                    string      `json:"Costcenter"`
	CostcenterDescription         string      `json:"CostcenterDescription"`
	Costunit                      string      `json:"Costunit"`
	CostunitDescription           string      `json:"CostunitDescription"`
	Created                       *types.Date `json:"Created"`
	Creator                       types.GUID  `json:"Creator"`
	CreatorFullName               string      `json:"CreatorFullName"`
	DeductionPercentage           float64     `json:"DeductionPercentage"`
	DepreciatedAmount             float64     `json:"DepreciatedAmount"`
	DepreciatedPeriods            int32       `json:"DepreciatedPeriods"`
	DepreciatedStartDate          *types.Date `json:"DepreciatedStartDate"`
	Description                   string      `json:"Description"`
	Division                      int32       `json:"Division"`
	EndDate                       *types.Date `json:"EndDate"`
	EngineEmission                int16       `json:"EngineEmission"`
	EngineType                    int16       `json:"EngineType"`
	GLTransactionLine             types.GUID  `json:"GLTransactionLine"`
	GLTransactionLineDescription  string      `json:"GLTransactionLineDescription"`
	InvestmentAccount             types.GUID  `json:"InvestmentAccount"`
	InvestmentAccountCode         string      `json:"InvestmentAccountCode"`
	InvestmentAccountName         string      `json:"InvestmentAccountName"`
	InvestmentAmountDC            float64     `json:"InvestmentAmountDC"`
	InvestmentAmountFC            float64     `json:"InvestmentAmountFC"`
	InvestmentCurrency            string      `json:"InvestmentCurrency"`
	InvestmentCurrencyDescription string      `json:"InvestmentCurrencyDescription"`
	InvestmentDate                *types.Date `json:"InvestmentDate"`
	InvestmentDeduction           int16       `json:"InvestmentDeduction"`
	Modified                      *types.Date `json:"Modified"`
	Modifier                      types.GUID  `json:"Modifier"`
	ModifierFullName              string      `json:"ModifierFullName"`
	Notes                         string      `json:"Notes"`
	Parent                        types.GUID  `json:"Parent"`
	ParentCode                    string      `json:"ParentCode"`
	ParentDescription             string      `json:"ParentDescription"`
	//Picture                       `json:"Picture"`
	PictureFileName          string      `json:"PictureFileName"`
	PrimaryMethod            types.GUID  `json:"PrimaryMethod"`
	PrimaryMethodCode        string      `json:"PrimaryMethodCode"`
	PrimaryMethodDescription string      `json:"PrimaryMethodDescription"`
	ResidualValue            float64     `json:"ResidualValue"`
	StartDate                *types.Date `json:"StartDate"`
	Status                   int16       `json:"Status"`
	TransactionEntryID       types.GUID  `json:"TransactionEntryID"`
	TransactionEntryNo       int32       `json:"TransactionEntryNo"`
}

func (eo *ExactOnline) GetAssetsInternal(filter string) (*[]Asset, error) {
	selectFields := GetJsonTaggedFieldNames(Asset{})
	urlStr := fmt.Sprintf("%s%s/assets/Assets?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	assets := []Asset{}

	for urlStr != "" {
		its := []Asset{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetAssetsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		assets = append(assets, its...)

		urlStr = str
		//urlStr = ""
	}

	return &assets, nil
}

func (eo *ExactOnline) GetAssets() (*[]Asset, error) {
	acc, err := eo.GetAssetsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}

/*
func (eo *ExactOnline) GetAssets() error {
	selectFields := GetJsonTaggedFieldNames(Asset{})
	urlStr := fmt.Sprintf("%s%s/logistics/Assets?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	//fmt.Println(urlStr)

	for urlStr != "" {
		it := []Asset{}

		str, err := eo.Get(urlStr, &it)
		if err != nil {
			fmt.Println("ERROR in GetAssets:", err)
			fmt.Println("url:", urlStr)
			return err
		}

		eo.Assets = append(eo.Assets, it...)

		urlStr = str
	}

	return nil
}*/
