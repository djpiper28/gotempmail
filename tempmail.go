package gotempmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// The TempMail struct stores the state of the tempmail
// along with error information
type TempMail struct {
	email    string
	password string
	id       string
	jwt      string
	// Any errors whilst building the TempMail are stored here
	Err error
}

// Inits a new TempMail instance, this is part of a builder type constructor,
// You should also .Address().Password().CreateAccount() to get a usable object
func New() *TempMail {
	ret := TempMail{}
	return &ret
}

// Sets the email address
func (tm *TempMail) Address(address string) *TempMail {
	tm.email = address
	return tm
}

// Sets the password
func (tm *TempMail) Password(password string) *TempMail {
	tm.password = password
	return tm
}

// Validates that the account can be made
func (tm *TempMail) Validate() error {
	if len(tm.password) == 0 {
		return fmt.Errorf("NO PASSWORD")
	}

	if len(tm.email) == 0 {
		return fmt.Errorf("NO EMAIL")
	}

	if len(tm.email) < 3 {
		return fmt.Errorf("EMAIL '%s' CANNOT BE VALID", tm.email)
	}

	if !strings.Contains(tm.email, "@") {
		return fmt.Errorf("NO @ IN EMAIL '%s'", tm.email)
	}
	return nil
}

type createAccountJson struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type createAccountRespJson struct {
	/* Other information that is not important
	   @context: /contexts/Account
	   @id: /accounts/6478c277752952794a10c466
	   @type: Account
	   address: testing1685635703@internetkeno.com
	   createdAt: 2023-06-01T16:08:23+00:00
	   quota: 40000000
	   updatedAt: 2023-06-01T16:08:23+00:00
	   used: 0
	*/
	Id       string `json:"id"`
	Deleted  bool   `json:"isDeleted"`
	Disabled bool   `json:"isDisabled"`
}

func (tm *TempMail) createAccount() error {
	err := tm.Validate()
	if err != nil {
		return fmt.Errorf("VALIDATION ERROR %s", err)
	}

	tmp := createAccountJson{Address: tm.email,
		Password: tm.password}
	msgBody, err := json.Marshal(tmp)
	if err != nil {
		return err
	}

	resp, err := http.Post(ACCOUNT_REGISTER_LINK,
		JSON_CONTENT,
		bytes.NewBuffer(msgBody))
	if err != nil {
		return fmt.Errorf("CANNOT POST %s", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return StatusCodeErr(resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BodyReadErr(err)
	}

	var respBody createAccountRespJson
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return JsonParseErr(err)
	}

	if respBody.Deleted {
		return fmt.Errorf("EMAIL DELETED")
	}

	if respBody.Disabled {
		return fmt.Errorf("EMAIL DISABLED")
	}

	tm.id = respBody.Id
	return nil
}

type authReqJson struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type authRespJson struct {
	Token string `json:"token"`
}

// Refreshes the authentication for the account, usually this does not get called
func (tm *TempMail) RefreshAuth() error {
	tmp := authReqJson{Address: tm.email,
		Password: tm.password}
	msgBody, err := json.Marshal(&tmp)
	if err != nil {
		return err
	}

	resp, err := http.Post(AUTH_LINK,
		JSON_CONTENT,
		bytes.NewBuffer(msgBody))

	if err != nil {
		return fmt.Errorf("CANNOT POST %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return StatusCodeErr(resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BodyReadErr(err)
	}

	var respBody authRespJson
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return JsonParseErr(err)
	}

	if respBody.Token == "" {
		return fmt.Errorf("EMPTY AUTH TOKEN")
	}
	tm.jwt = respBody.Token

	return nil
}

// Creates the account on the TempMail server, this is the last bit of the builder functions
func (tm *TempMail) CreateAccount() *TempMail {
	tm.Err = tm.createAccount()
	// Fail fast
	if tm.Err != nil {
		return tm
	}

	tm.Err = tm.RefreshAuth()
	return tm
}

// Gets the login data for the request headers
func (tm *TempMail) getLoginData() string {
	return "Bearer " + tm.jwt
}

type EmailAddr struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// An email in the inbox
type Email struct {
	// Sender email address (from:)
	Sender EmailAddr `json:"from"`
	// People who recieved the email (to:)
	Receipient []EmailAddr `json:"to"`
	// Subject line
	Subject string `json:"subject"`
	// The first bit of the body
	Intro          string `json:"intro"`
	HasAttachments bool   `json:"hasAttachments"`
	Size           int    `json:"size"`
	Seen           bool   `json:"seen"`
	downloadUrl    string `json:"downloadUrl"`
	id             string `json:"id"`
}

type emailsJson struct {
	Emails []Email `json:"hydra:member"`
}

// Gets the emails for an TempMail object
func (tm *TempMail) GetEmails() ([]Email, error) {
	req, err := http.NewRequest(http.MethodGet, MESSAGES_LINK, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(AUTH_HEADER, tm.getLoginData())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CANNOT GET MESSAGES %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, StatusCodeErr(resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, BodyReadErr(err)
	}

	var emails emailsJson
	err = json.Unmarshal(body, &emails)
	if err != nil {
		return nil, JsonParseErr(err)
	}

	return emails.Emails, nil
}
