package exactonline

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// Contact stores Contact from exactonline
//
type Contact struct {
	ID            types.GUID `json:"ID"`
	Account       types.GUID `json:"Account"`
	Initials      string     `json:"Initials"`
	BirthName     string     `json:"BirthName"`
	FirstName     string     `json:"FirstName"`
	LastName      string     `json:"LastName"`
	Gender        string     `json:"Gender"`
	Title         string     `json:"Title"`
	Email         string     `json:"Email"`
	IsMainContact bool       `json:"IsMainContact"`
}

/*
var oldContact *Contact

// SaveValues saves current values in local copy of Contact
//
func (c *Contact) SaveValues(inserted bool) {
	oldContact = nil
	if !inserted {
		oldContact = new(Contact)
		oldContact.Initials = c.Initials
		oldContact.BirthName = c.BirthName
		oldContact.FirstName = c.FirstName
		oldContact.LastName = c.LastName
		oldContact.Gender = c.Gender
		oldContact.Title = c.Title
		oldContact.Email = c.Email
	}
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
}*/
/*
func (c *Contact) Values() (string, string) {
	old := ""
	new := ""

	if oldContact != nil {
		oldContact.Initials = strings.Trim(oldContact.Initials, " ")
		oldContact.BirthName = strings.Trim(oldContact.BirthName, " ")
		oldContact.FirstName = strings.Trim(oldContact.FirstName, " ")
		oldContact.LastName = strings.Trim(oldContact.LastName, " ")
		oldContact.Gender = strings.Trim(oldContact.Gender, " ")
		oldContact.Title = strings.Trim(oldContact.Title, " ")
	}

	c.Initials = strings.Trim(c.Initials, " ")
	c.BirthName = strings.Trim(c.BirthName, " ")
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
		new += ",BirthName:" + c.BirthName
	} else if oldContact.BirthName != c.BirthName {
		old += ",BirthName:" + oldContact.BirthName
		new += ",BirthName:" + c.BirthName
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
}*/

func (eo *ExactOnline) GetContactsInternal(filter string) (*[]Contact, *errortools.Error) {
	selectFields := utilities.GetTaggedFieldNames("json", Contact{})
	urlStr := fmt.Sprintf("%s/crm/Contacts?$select=%s", eo.baseURL(), selectFields)
	if filter != "" {
		urlStr += fmt.Sprintf("&$filter=%s", filter)
	}
	//fmt.Println(urlStr)

	contacts := []Contact{}

	for urlStr != "" {
		co := []Contact{}

		str, e := eo.Get(urlStr, &co)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, co...)

		urlStr = str
	}

	return &contacts, nil
}

func (eo *ExactOnline) GetContacts() *errortools.Error {
	co, e := eo.GetContactsInternal("")
	if e != nil {
		return e
	}
	eo.Contacts = *co

	return nil
}

func (eo *ExactOnline) GetContactsByEmail(account string, email string) (*[]Contact, *errortools.Error) {
	filter := fmt.Sprintf("Account eq guid'%s' and Email eq '%s'", account, email)

	contacts, e := eo.GetContactsInternal(filter)
	if e != nil {
		return nil, e
	}

	return contacts, nil
}

func (eo *ExactOnline) GetContactsByFullName(account string, fullname string) (*[]Contact, *errortools.Error) {
	filter := fmt.Sprintf("Account eq guid'%s' and FullName eq '%s'", account, fullname)

	contacts, e := eo.GetContactsInternal(filter)
	if e != nil {
		return nil, e
	}

	return contacts, nil
}

func (eo *ExactOnline) UpdateContact(c *Contact) *errortools.Error {
	urlStr := fmt.Sprintf("%s/crm/Contacts(guid'%s')", eo.baseURL(), c.ID.String())

	data := make(map[string]string)
	data["Initials"] = c.Initials
	data["BirthName"] = c.BirthName
	data["FirstName"] = c.FirstName
	data["LastName"] = c.LastName
	data["Gender"] = c.Gender
	data["Title"] = c.Title
	data["Email"] = c.Email

	e := eo.PutValues(urlStr, data)
	if e != nil {
		return e
	}

	return nil
}

func (eo *ExactOnline) InsertContact(c *Contact) *errortools.Error {
	urlStr := fmt.Sprintf("%s/crm/Contacts", eo.baseURL())

	data := make(map[string]string)
	data["Account"] = c.Account.String()
	data["Initials"] = c.Initials
	data["BirthName"] = c.BirthName
	data["FirstName"] = c.FirstName
	data["LastName"] = c.LastName
	data["Gender"] = c.Gender
	data["Title"] = c.Title
	data["Email"] = c.Email

	co := Contact{}

	e := eo.PostValues(urlStr, data, &co)
	if e != nil {
		return e
	}

	c.ID = co.ID

	return nil
}
