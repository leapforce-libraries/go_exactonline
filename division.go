package exactonline

import (
	"fmt"
	"strconv"
)

// Division stores division from exactonline
//
type Division struct {
	Code int `json:"Code"`
}

func (eo *ExactOnline) getDivisions() error {
	urlStr := fmt.Sprintf("%s%s/hrm/Divisions", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision))

	eo.Divisions = []Division{}

	for urlStr != "" {
		d := []Division{}

		str, err := eo.get(urlStr, &d)
		if err != nil {
			return err
		}

		eo.Divisions = append(eo.Divisions, d...)

		urlStr = str
		//urlStr = "" //temp
	}

	return nil
}
