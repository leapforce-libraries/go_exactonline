package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Division stores division from exactonline
//
type Division struct {
	Code int `json:"Code"`
}

func (eo *ExactOnline) GetDivisions() *errortools.Error {
	urlStr := fmt.Sprintf("%s/hrm/Divisions", eo.baseURL())

	eo.Divisions = []Division{}

	for urlStr != "" {
		d := []Division{}

		str, err := eo.Get(urlStr, &d)
		if err != nil {
			return err
		}

		eo.Divisions = append(eo.Divisions, d...)

		urlStr = str
		//urlStr = "" //temp
	}

	return nil
}
