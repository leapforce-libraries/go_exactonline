package exactonline

import (
	"encoding/json"
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
)

// Me stores Me from exactonline
//
type Me struct {
	UserID                      types.GUID      `json:"UserID"`
	AccountingDivision          int32           `json:"AccountingDivision"`
	CurrentDivision             int32           `json:"CurrentDivision"`
	DivisionCustomer            types.GUID      `json:"DivisionCustomer"`
	DivisionCustomerCode        string          `json:"DivisionCustomerCode"`
	DivisionCustomerName        string          `json:"DivisionCustomerName"`
	DivisionCustomerSiretNumber string          `json:"DivisionCustomerSiretNumber"`
	DivisionCustomerVatNumber   string          `json:"DivisionCustomerVatNumber"`
	DossierDivision             int32           `json:"DossierDivision"`
	Email                       string          `json:"Email"`
	EmployeeID                  types.GUID      `json:"EmployeeID"`
	FirstName                   string          `json:"FirstName"`
	FullName                    string          `json:"FullName"`
	Gender                      string          `json:"Gender"`
	Initials                    string          `json:"Initials"`
	IsClientUser                bool            `json:"IsClientUser"`
	IsMyFirmPortalUser          bool            `json:"IsMyFirmPortalUser"`
	Language                    string          `json:"Language"`
	LanguageCode                string          `json:"LanguageCode"`
	LastName                    string          `json:"LastName"`
	Legislation                 string          `json:"Legislation"`
	MiddleName                  string          `json:"MiddleName"`
	Mobile                      string          `json:"Mobile"`
	Nationality                 string          `json:"Nationality"`
	Phone                       string          `json:"Phone"`
	PhoneExtension              string          `json:"PhoneExtension"`
	PictureUrl                  string          `json:"PictureUrl"`
	ServerTime                  string          `json:"ServerTime"`
	ServerUtcOffset             float64         `json:"ServerUtcOffset"`
	ThumbnailPicture            json.RawMessage `json:"ThumbnailPicture"`
	ThumbnailPictureFormat      string          `json:"ThumbnailPictureFormat"`
	Title                       string          `json:"Title"`
	UserName                    string          `json:"UserName"`
}

func (eo *ExactOnline) GetMe() (*Me, *errortools.Error) {
	urlStr := fmt.Sprintf("%s/current/Me", eo.apiURL())

	me := []Me{}

	_, err := eo.Get(urlStr, &me)
	if err != nil {
		return nil, err
	}

	return &me[0], nil
}
