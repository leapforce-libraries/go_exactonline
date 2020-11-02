package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Bank stores Bank from exactonline
//
type Bank struct {
	ID              types.GUID  `json:"ID"`
	BankName        string      `json:"BankName"`
	BICCode         string      `json:"BICCode"`
	Country         string      `json:"Country"`
	Created         *types.Date `json:"Created"`
	Description     string      `json:"Description"`
	Format          string      `json:"Format"`
	HomePageAddress string      `json:"HomePageAddress"`
	Modified        *types.Date `json:"Modified"`
	Status          string      `json:"Status"`
}

func (eo *ExactOnline) GetBanksInternal(filter string) (*[]Bank, error) {
	selectFields := utilities.GetTaggedFieldNames("json", Bank{})
	urlStr := fmt.Sprintf("%s/cashflow/Banks?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	banks := []Bank{}

	for urlStr != "" {
		its := []Bank{}

		str, err := eo.Get(urlStr, &its)
		if err != nil {
			fmt.Println("ERROR in GetBanksInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		banks = append(banks, its...)

		urlStr = str
		//urlStr = ""
	}

	return &banks, nil
}

func (eo *ExactOnline) GetBanks() (*[]Bank, error) {
	acc, err := eo.GetBanksInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
