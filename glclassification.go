package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
)

// GLClassification stores GLClassification from exactonline
//
type GLClassification struct {
	ID                           types.GUID  `json:"ID"`
	Abstract                     bool        `json:"Abstract"`
	Balance                      string      `json:"Balance"`
	Code                         string      `json:"Code"`
	Created                      *types.Date `json:"Created"`
	Creator                      types.GUID  `json:"Creator"`
	CreatorFullName              string      `json:"CreatorFullName"`
	Description                  string      `json:"Description"`
	Division                     int32       `json:"Division"`
	IsTupleSubElement            bool        `json:"IsTupleSubElement"`
	Modified                     *types.Date `json:"Modified"`
	Modifier                     types.GUID  `json:"Modifier"`
	ModifierFullName             string      `json:"ModifierFullName"`
	Name                         string      `json:"Name"`
	Nillable                     bool        `json:"Nillable"`
	Parent                       types.GUID  `json:"Parent"`
	PeriodType                   string      `json:"PeriodType"`
	SubstitutionGroup            string      `json:"SubstitutionGroup"`
	TaxonomyNamespace            types.GUID  `json:"TaxonomyNamespace"`
	TaxonomyNamespaceDescription string      `json:"TaxonomyNamespaceDescription"`
	Type                         types.GUID  `json:"Type"`
}

func (eo *ExactOnline) GetGLClassificationsInternal(filter string) (*[]GLClassification, error) {
	selectFields := GetJsonTaggedFieldNames(GLClassification{})
	urlStr := fmt.Sprintf("%s%s/bulk/financial/GLClassifications?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	glClassifications := []GLClassification{}

	for urlStr != "" {
		ac := []GLClassification{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetGLClassificationsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		glClassifications = append(glClassifications, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &glClassifications, nil
}

func (eo *ExactOnline) GetGLClassifications() (*[]GLClassification, error) {
	acc, err := eo.GetGLClassificationsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
