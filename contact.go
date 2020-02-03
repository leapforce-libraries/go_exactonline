package exactonline

import (
	"fmt"
	"strconv"
	"strings"

	types "github.com/Leapforce-nl/go_types"
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

var oldContact *Contact

// SaveValues saves current values in local copy of Contact
//
func (c *Contact) SaveValues() {
	oldContact = new(Contact)
	oldContact.Initials = c.Initials
	oldContact.FirstName = c.FirstName
	oldContact.LastName = c.LastName
	oldContact.Gender = c.Gender
	oldContact.Title = c.Title
	oldContact.Email = c.Email
}

// Values return comma separated values of Account
//
/*func (c *Contact) Values() string {
	return fmt.Sprintf("Initials: %s, FirstName: %s, LastName: %s, Gender: %s, Title: %s, Email: %s",
		c.Initials,
		c.FirstName,
		c.LastName,
		c.Gender,
		c.Title,
		c.Email)
}*/

func (c *Contact) Values() (string, string) {
	old := ""
	new := ""

	if oldContact != nil {
		oldContact.Initials = strings.Trim(oldContact.Initials, " ")
		oldContact.FirstName = strings.Trim(oldContact.FirstName, " ")
		oldContact.LastName = strings.Trim(oldContact.LastName, " ")
		oldContact.Gender = strings.Trim(oldContact.Gender, " ")
		oldContact.Title = strings.Trim(oldContact.Title, " ")
	}

	c.Initials = strings.Trim(c.Initials, " ")
	c.FirstName = strings.Trim(c.FirstName, " ")
	c.LastName = strings.Trim(c.LastName, " ")
	c.Gender = strings.Trim(c.Gender, " ")
	c.Title = strings.Trim(c.Title, " ")

	if oldContact == nil {
		new += ",Initials:" + c.Initials
	} else if oldContact.Initials != c.Initials {
		old += ",Initials:" + oldContact.Initials
		new += ",Initials:" + c.Initials
	}

	if oldContact == nil {
		new += ",FirstName:" + c.FirstName
	} else if oldContact.FirstName != c.FirstName {
		old += ",FirstName:" + oldContact.FirstName
		new += ",FirstName:" + c.FirstName
	}

	if oldContact == nil {
		new += ",LastName:" + c.LastName
	} else if oldContact.LastName != c.LastName {
		old += ",LastName:" + oldContact.LastName
		new += ",LastName:" + c.LastName
	}

	if oldContact == nil {
		new += ",Gender:" + c.Gender
	} else if oldContact.Gender != c.Gender {
		old += ",Gender:" + oldContact.Gender
		new += ",Gender:" + c.Gender
	}

	if oldContact == nil {
		new += ",Title:" + c.Title
	} else if oldContact.Title != c.Title {
		old += ",Title:" + oldContact.Title
		new += ",Title:" + c.Title
	}

	if oldContact == nil {
		new += ",Email:" + c.Email
	} else if oldContact.Email != c.Email {
		old += ",Email:" + oldContact.Email
		new += ",Email:" + c.Email
	}

	old = strings.TrimLeft(old, ",")
	new = strings.TrimLeft(new, ",")

	return old, new
}

func (eo *ExactOnline) GetContactsInternal(filter string) (*[]Contact, error) {
	selectFields := GetJsonTaggedFieldNames(Contact{})
	urlStr := fmt.Sprintf("%s%s/crm/Contacts?$select=%s", eo.ApiUrl, strconv.Itoa(eo.Division), selectFields)
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
	urlStr := fmt.Sprintf("%s%s/crm/Contacts(guid'%s')", eo.ApiUrl, strconv.Itoa(eo.Division), c.ID.String())

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
	urlStr := fmt.Sprintf("%s%s/crm/Contacts", eo.ApiUrl, strconv.Itoa(eo.Division))

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
