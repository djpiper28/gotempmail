package gotempmail

import (
	"net/http"
)

type TempMail struct {
	email string
}

const (
	// The base URL of the Temp Mail service, this might change tbh
	BASE_URL            = "https://api.temp-mail.ru"
	DELETE_MESSAGE_LINK = BASE_URL + "/request/delete/id/md5/"
	VIEW_MESSAGE_LINK   = BASE_URL + "/request/mail/id/md5/"
	DOMAIN_LIST_LINK    = BASE_URL + "/request/domains"
)

// Gets all of the TempMail domains
func GetDomains() ([]string, error) {
	_, err := http.Get(DOMAIN_LIST_LINK)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Inits a new TempMail instance
func New() *TempMail {
	ret := TempMail{}
	return &ret
}
