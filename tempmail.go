package gotempmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// The TempMail struct stores the state of the tempmail
// along with error information
type TempMail struct {
	email    string
	password string
	// Any errors whilst building the TempMail are stored here
	Err error
}

// Inits a new TempMail instance
func New() *TempMail {
	ret := TempMail{}
	return &ret
}

// Sets the email address
func (tm *TempMail) Address(address string) *TempMail {
	tm.email = address
	return tm
}

func (tm *TempMail) Password(password string) *TempMail {
	tm.password = password
	return tm
}

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

	log.Print(string(body))

	return fmt.Errorf("TODO")
}

// Creates the account on the TempMail server
func (tm *TempMail) CreateAccount() *TempMail {
	tm.Err = tm.createAccount()
	return tm
}
