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

type domainJson struct {
	Domain string `json:"domain"`
	/* Other information that is not needed
	@id: /domains/64637851672bde8f395a0b1a
	@type: Domain
	createdAt: 2023-05-16T00:00:00+00:00
	domain: internetkeno.com
	id: 64637851672bde8f395a0b1a
	isActive: true
	isPrivate: false
	updatedAt: 2023-05-16T00:00:00+00:00
	*/
}

type domainsJson struct {
	/* Other information that is not needed
	@context: /contexts/Domain
	@id: /domains
	@type: hydra:Collection
	*/
	Domains []domainJson `json:"hydra:member"`
}

const (
	// Content type
	JSON_CONTENT = "application/json"
	// The base URL of the Temp Mail service, this might change tbh
	BASE_URL              = "https://api.mail.tm"
	DOMAIN_LIST_LINK      = BASE_URL + "/domains"
	ACCOUNT_REGISTER_LINK = BASE_URL + "/accounts"
)

// Gets all of the TempMail domains
func GetDomains() ([]string, error) {
	resp, err := http.Get(DOMAIN_LIST_LINK)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UNEXPECTED RETURN CODE (%d)",
			resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CANNOT READ BODY %s", err)
	}

	var domains domainsJson
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, fmt.Errorf("CANNOT PARSE DOMAINS %s", err)
	}

	ret := make([]string, len(domains.Domains))
	for i, domain := range domains.Domains {
		ret[i] = domain.Domain
	}

	return ret, nil
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
