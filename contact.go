package exactonline

import (
	"fmt"
	"strconv"

	types "github.com/leapforce-nl/go_types"
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

func (eo *ExactOnline) GetContactsInternal(filter string) (*[]Contact, error) {
	selectFields := GetJsonTaggedFieldNames(Contact{})
	urlStr := fmt.Sprintf("%s%s/crm/Contacts?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Me.CurrentDivision), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	contacts := []Contact{}

	for urlStr != "" {
		co := []Contact{}

		str, err := eo.Get(urlStr, &co)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, co...)

		urlStr = str
	}

	return &contacts, nil
}

func (eo *ExactOnline) GetContacts() error {
	co, err := eo.GetContactsInternal("")
	if err != nil {
		return err
	}
	eo.Contacts = *co

	return nil
}

func (eo *ExactOnline) GetContactsByEmail(email string) ([]Contact, error) {
	filter := fmt.Sprintf("Email eq '%s'", email)
	contacts := []Contact{}

	co, err := eo.GetContactsInternal(filter)
	if err != nil {
		return contacts, nil
	}
	contacts = *co

	//fmt.Println("GetContactsByEmail:", email, "len:", len(contacts))

	return contacts, nil
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

	err := eo.Put(urlStr, data)
	if err != nil {
		return err
	}

	fmt.Println("\nUPDATED Contact")
	fmt.Println("url:", urlStr)
	fmt.Println("data:", data)

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

	err := eo.Post(urlStr, data, &co)
	if err != nil {
		fmt.Println("ERROR in InsertContact:", err)
		fmt.Println("url:", urlStr)
		fmt.Println("data:", data)
		return err
	}

	fmt.Println("\nINSERTED Contact", co.ID)
	fmt.Println("url:", urlStr)
	fmt.Println("data:", data)
	c.ID = co.ID

	return nil
}
