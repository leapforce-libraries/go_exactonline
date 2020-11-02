package exactonline

import (
	"fmt"

	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// GLScheme stores GLScheme from exactonline
//
type GLScheme struct {
	ID               types.GUID  `json:"ID"`
	Code             string      `json:"Code"`
	Created          *types.Date `json:"Created"`
	Creator          types.GUID  `json:"Creator"`
	CreatorFullName  string      `json:"CreatorFullName"`
	Description      string      `json:"Description"`
	Division         int32       `json:"Division"`
	Main             byte        `json:"Main"`
	Modified         *types.Date `json:"Modified"`
	Modifier         types.GUID  `json:"Modifier"`
	ModifierFullName string      `json:"ModifierFullName"`
	TargetNamespace  string      `json:"TargetNamespace"`
}

func (eo *ExactOnline) GetGLSchemesInternal(filter string) (*[]GLScheme, error) {
	selectFields := utilities.GetTaggedFieldNames("json", GLScheme{})
	urlStr := fmt.Sprintf("%s/financial/GLSchemes?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	glSchemes := []GLScheme{}

	for urlStr != "" {
		ac := []GLScheme{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetGLSchemesInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		glSchemes = append(glSchemes, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &glSchemes, nil
}

func (eo *ExactOnline) GetGLSchemes() (*[]GLScheme, error) {
	acc, err := eo.GetGLSchemesInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
