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
	"sync"
	"time"

	types "github.com/Leapforce-nl/go_types"
	"github.com/getsentry/sentry-go"
)

var tokenMutex sync.Mutex

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	Expiry       time.Time
	RefreshToken string `json:"refresh_token"`
}

type ApiError struct {
	Error       string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

func LockToken() {
	tokenMutex.Lock()
}

func UnlockToken() {
	tokenMutex.Unlock()
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
	if t.Expiry.Add(-60 * time.Second).Before(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (eo *ExactOnline) GetToken(data url.Values) error {
	guid := types.NewGUID()
	fmt.Println("GetTokenGUID:", guid)

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

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		fmt.Println("GetTokenGUID:", guid)
		fmt.Println("AccessToken:", eo.Token.AccessToken)
		fmt.Println("Refresh:", eo.Token.RefreshToken)
		fmt.Println("Expiry:", eo.Token.Expiry)
		fmt.Println("Now:", time.Now())

		eoError := ApiError{}

		err = json.Unmarshal(b, &eoError)
		message := ""
		if err == nil {
			message = fmt.Sprintln("Error:", eoError.Error, ", ", eoError.Description)
			fmt.Println(message)
			//return err
		}

		if res.StatusCode == 401 {
			if eo.IsLive {
				sentry.CaptureMessage("ExactOnline refreshtoken not valid, login needed to retrieve a new one. Error: " + message)
			}
			eo.InitToken()
		}

		return &types.ErrorString{fmt.Sprintf("Server returned statuscode %v, url: %s", res.StatusCode, req.URL)}
	}

	token := Token{}

	err = json.Unmarshal(b, &token)
	if err != nil {
		return err
	}

	fmt.Println("old token:")
	fmt.Println(eo.Token.AccessToken)
	fmt.Println("old refresh token:")
	fmt.Println(eo.Token.RefreshToken)
	fmt.Println("old expiry:")
	fmt.Println(eo.Token.Expiry)

	expiresIn, err := strconv.ParseInt(token.ExpiresIn, 10, 64)
	if err != nil {
		return err
	}
	token.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	//token.Expiry = time.Now().Add(time.Duration(60) * time.Second)
	//eo.Token = &token
	eo.Token.Expiry = token.Expiry
	eo.Token.RefreshToken = token.RefreshToken
	eo.Token.AccessToken = token.AccessToken

	err = eo.SaveTokenToBigQuery()
	if err != nil {
		return err
	}

	fmt.Println("new token:")
	fmt.Println(eo.Token.AccessToken)
	fmt.Println("new refresh token:")
	fmt.Println(eo.Token.RefreshToken)
	fmt.Println("new expiry:")
	fmt.Println(eo.Token.Expiry)
	fmt.Println("GetTokenGUID:", guid)

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
	fmt.Println("***GetTokenFromRefreshToken***")
	//fmt.Println(eo.Token.RefreshToken[0:20])

	//always get refresh token from BQ prior to using it
	eo.GetTokenFromBigQuery()

	if !eo.Token.Refreshable() {
		err := eo.InitToken()
		if err != nil {
			return err
		}
		//return &types.ErrorString{"Token is not valid."}
	}
	data := url.Values{}
	data.Set("client_id", eo.ClientID)
	data.Set("client_secret", eo.ClientSecret)
	data.Set("refresh_token", eo.Token.RefreshToken)
	data.Set("grant_type", "refresh_token")

	return eo.GetToken(data)
}

func (eo *ExactOnline) ValidateToken() error {
	LockToken()
	defer UnlockToken()

	if !eo.Token.Useable() {
		err := eo.GetTokenFromRefreshToken()
		if err != nil {
			return err
		}

		if !eo.Token.Useable() {
			if eo.IsLive {
				sentry.CaptureMessage("ExactOnline refreshtoken not found or empty, login needed to retrieve a new one.")
			}
			err := eo.InitToken()
			if err != nil {
				return err
			}
			//return &types.ErrorString{""}
		}
	}

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
