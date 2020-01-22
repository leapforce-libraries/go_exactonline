package exactonline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	types "github.com/leapforce-nl/go_types"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	Expiry       time.Time
	RefreshToken string `json:"refresh_token"`
}

func (t *Token) Useable() bool {
	if t == nil {
		return false
	}
	if t.AccessToken == "" || t.RefreshToken == "" {
		return false
	}
	return true
}

func (t *Token) Refreshable() bool {
	if t == nil {
		return false
	}
	if t.RefreshToken == "" {
		return false
	}
	return true
}

func (t *Token) IsExpired() (bool, error) {
	if !t.Useable() {
		return true, &types.ErrorString{"Token is not valid."}
	}
	if t.Expiry.Add(-30 * time.Second).Before(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (eo *ExactOnline) GetToken(data url.Values) error {
	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodPost, eo.TokenURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		return err
	}

	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v", res.StatusCode)}
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	token := Token{}

	err = json.Unmarshal(b, &token)
	if err != nil {
		return err
	}

	expiresIn, err := strconv.ParseInt(token.ExpiresIn, 10, 64)
	if err != nil {
		return err
	}
	token.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	eo.Token = &token

	err = eo.SaveTokenToBigQuery()
	if err != nil {
		return err
	}

	fmt.Println("new token:")
	fmt.Println(eo.Token.AccessToken)

	return nil
}

func (eo *ExactOnline) GetTokenFromCode(code string) error {
	//fmt.Println("GetTokenFromCode")
	data := url.Values{}
	data.Set("client_id", eo.ClientID)
	data.Set("client_secret", eo.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", eo.RedirectURL)

	return eo.GetToken(data)
}

func (eo *ExactOnline) GetTokenFromRefreshToken() error {
	//fmt.Println("GetTokenFromRefreshToken")
	//fmt.Println(eo.Token.RefreshToken[0:20])
	if !eo.Token.Refreshable() {
		return &types.ErrorString{"Token is not valid."}
	}
	data := url.Values{}
	data.Set("client_id", eo.ClientID)
	data.Set("client_secret", eo.ClientSecret)
	data.Set("refresh_token", eo.Token.RefreshToken)
	data.Set("grant_type", "refresh_token")

	return eo.GetToken(data)
}

func (eo *ExactOnline) ValidateToken() error {
	if !eo.Token.Useable() {
		if !eo.Token.Refreshable() {
			err := eo.GetTokenFromBigQuery()
			if err != nil {
				return err
			}
			//fmt.Println(time.Now(), eo.Token.Expiry, "[from bq]", eo.Token.AccessToken, eo.Token.RefreshToken)
		}

		if eo.Token.Refreshable() {
			err := eo.GetTokenFromRefreshToken()
			if err != nil {
				return err
			}
		}

		if !eo.Token.Useable() {
			err := eo.InitToken()
			if err != nil {
				return err
			}
			//return &types.ErrorString{""}
		}
	}

	//fmt.Println("[try]", time.Now(), eo.Token.Expiry, "[me]", eo.Me.CurrentDivision, "[at]", eo.Token.AccessToken[0:20], "[rt]", eo.Token.RefreshToken[0:20])

	isExpired, err := eo.Token.IsExpired()
	if err != nil {
		return err
	}
	if isExpired {
		//fmt.Println(time.Now(), "[token expired]")
		err = eo.GetTokenFromRefreshToken()
		if err != nil {
			return err
		}
	}

	//fmt.Println("[done]", time.Now(), eo.Token.Expiry, "[me]", eo.Me.CurrentDivision, "[at]", eo.Token.AccessToken[0:20], "[rt]", eo.Token.RefreshToken[0:20])

	return nil
}

func (eo *ExactOnline) InitToken() error {
	if eo == nil {
		return &types.ErrorString{"ExactOnline variable not initialized"}
	}

	url := fmt.Sprintf("https://start.exactonline.nl/api/oauth2/auth?client_id=%s&response_type=code&redirect_uri=%s", eo.ClientID, eo.RedirectURL)

	fmt.Println("Go to this url to get new access token:\n")
	fmt.Println(url + "\n")

	// Create a new redirect route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		//
		// get authorization code
		//
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		err = eo.GetTokenFromCode(code)
		if err != nil {
			fmt.Println(err)
		}

		w.WriteHeader(http.StatusFound)

		return
	})

	http.ListenAndServe(":8080", nil)

	return nil
}
