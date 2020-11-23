package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// GLAccountClassificationMapping stores GLAccountClassificationMapping from exactonline
//
type GLAccountClassificationMapping struct {
	ID                        types.GUID `json:"ID"`
	Classification            types.GUID `json:"Classification"`
	ClassificationCode        string     `json:"ClassificationCode"`
	ClassificationDescription string     `json:"ClassificationDescription"`
	Division                  int64      `json:"Division"`
	GLAccount                 types.GUID `json:"GLAccount"`
	GLAccountCode             string     `json:"GLAccountCode"`
	GLAccountDescription      string     `json:"GLAccountDescription"`
	GLSchemeCode              string     `json:"GLSchemeCode"`
	GLSchemeDescription       string     `json:"GLSchemeDescription"`
	GLSchemeID                types.GUID `json:"GLSchemeID"`
}

func (eo *ExactOnline) GetGLAccountClassificationMappingsInternal(filter string) (*[]GLAccountClassificationMapping, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", GLAccountClassificationMapping{})
	urlStr := fmt.Sprintf("%s/financial/GLAccountClassificationMappings?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	glAccountClassificationMappings := []GLAccountClassificationMapping{}

	for urlStr != "" {
		ac := []GLAccountClassificationMapping{}

		str, err := eo.Get(urlStr, &ac)
		if err != nil {
			fmt.Println("ERROR in GetGLAccountClassificationMappingsInternal:", err)
			fmt.Println("url:", urlStr)
			return nil, err
		}

		glAccountClassificationMappings = append(glAccountClassificationMappings, ac...)

		urlStr = str
		//urlStr = ""
	}

	return &glAccountClassificationMappings, nil
}

func (eo *ExactOnline) GetGLAccountClassificationMappings() (*[]GLAccountClassificationMapping, *errortools.Error) {
	acc, err := eo.GetGLAccountClassificationMappingsInternal("")
	if err != nil {
		return nil, err
	}

	return acc, nil
}
