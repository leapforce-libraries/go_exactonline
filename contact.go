package exactonline

import (
	"fmt"
	"strconv"
	"types"
)

// Contact stores Contact from exactonline
//
type Contact struct {
	ID            types.GUID `json:"ID"`
	Account       types.GUID `json:"Account"`
	Initials      string     `json:"Initials"`
	FirstName     string     `json:"FirstName"`
	LastName      string     `json:"LastName"`
	Gender        string     `json:"Gender"`
	Title         string     `json:"Title"`
	Email         string     `json:"Email"`
	IsMainContact bool       `json:"IsMainContact"`
}

// Values return comma separated values of Account
//
func (c *Contact) Values() string {
	return fmt.Sprintf("Initials: %s, FirstName: %s, LastName: %s, Gender: %s, Title: %s, Email: %s",
		c.Initials,
		c.FirstName,
		c.LastName,
		c.Gender,
		c.Title,
		c.Email)
}

func (eo *ExactOnline) getContacts() error {
	selectFields := GetJsonTaggedFieldNames(Contact{})
	urlStr := fmt.Sprintf("%s%s/crm/Contacts?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
	//fmt.Println(urlStr)

	eo.Contacts = []Contact{}

	for urlStr != "" {
		co := []Contact{}

		str, err := eo.get(urlStr, &co)
		if err != nil {
			return err
		}

		eo.Contacts = append(eo.Contacts, co...)

		urlStr = str
		//urlStr = "" //temp
	}

	return nil
}

func (eo *ExactOnline) UpdateContact(c *Contact) error {
	urlStr := fmt.Sprintf("%s%s/crm/Contacts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), c.ID.String())

	data := make(map[string]string)
	data["Initials"] = c.Initials
	data["FirstName"] = c.FirstName
	data["LastName"] = c.LastName
	data["Gender"] = c.Gender
	data["Title"] = c.Title
	data["Email"] = c.Email

	fmt.Println("update", urlStr, c.Email)

	err := eo.put(urlStr, data)
	if err != nil {
		return err
	}

	//time.Sleep(1 * time.Second)

	return nil
}

func (eo *ExactOnline) InsertContact(c *Contact) error {
	urlStr := fmt.Sprintf("%s%s/crm/Contacts", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision))

	data := make(map[string]string)
	data["Account"] = c.Account.String()
	data["Initials"] = c.Initials
	data["FirstName"] = c.FirstName
	data["LastName"] = c.LastName
	data["Gender"] = c.Gender
	data["Title"] = c.Title
	data["Email"] = c.Email

	co := Contact{}

	fmt.Println("insert", urlStr, c.Account.String(), c.Email)

	err := eo.post(urlStr, data, &co)
	if err != nil {
		fmt.Println(err)
		return err
	}

	c.ID = co.ID

	//fmt.Println("Inserted:", a.ID.String())

	return nil
}
