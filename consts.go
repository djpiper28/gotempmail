package gotempmail

const (
	// Content type
	JSON_CONTENT = "application/json"
	// Auth header for requests
	AUTH_HEADER = "Authorization"
	// The base URL of the Temp Mail service, this might change tbh
	BASE_URL              = "https://api.mail.tm"
	DOMAIN_LIST_LINK      = BASE_URL + "/domains"
	ACCOUNT_REGISTER_LINK = BASE_URL + "/accounts"
	MESSAGES_LINK         = BASE_URL + "/messages"
	AUTH_LINK             = BASE_URL + "/token"
)
