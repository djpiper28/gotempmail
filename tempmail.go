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
		return fmt.Errorf("UNEXPECTED RETURN CODE (%d)",
			resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("CANNOT READ BODY %s", err)
	}

	var respBody createAccountRespJson
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return fmt.Errorf("CANNOT PARSE JSON %s", err)
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

// Creates the account on the TempMail server
func (tm *TempMail) CreateAccount() *TempMail {
	tm.Err = tm.createAccount()
	return tm
}
