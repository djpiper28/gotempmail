package gotempmail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
