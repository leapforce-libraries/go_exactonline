package exactonline

// Me stores Me from exactonline
//
type Me struct {
	CurrentDivision int    `json:"CurrentDivision"`
	FirstName       string `json:"FirstName"`
}

func (eo *ExactOnline) GetMe() error {
	urlStr := "https://start.exactonline.nl/api/v1/current/Me"

	me := []Me{}

	_, err := eo.Get(urlStr, &me)
	if err != nil {
		return err
	}

	//eo.Me = me[0]

	return nil
}
