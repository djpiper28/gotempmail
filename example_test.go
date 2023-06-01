package gotempmail

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// An example of how to use the library
func ExampleTest(t *testing.T) {
	// Get domains
	domains, err := GetDomains()
	if err != nil {
		t.Error(err)
	}
	if len(domains) == 0 {
		t.Error("No domains")
	}

	// Build the tempmail object and test for errors
	tempmail := New().
		Address("account" + fmt.Sprintf("%d",
			time.Now().Unix()) + "@" + domains[0]).
		Password("password123").
		CreateAccount()
	if tempmail.Err != nil {
		t.Errorf("tempmail err is %s", tempmail.Err)
	}

	// Imagine there was an email sent lmao
	// Get emails
	emails, err := tempmail.GetEmails()
	if err != nil {
		t.Errorf("tempmail GetEmails () err %s", err)
	}
	log.Printf("There are %d emails", len(emails))

	// Get email details
	for _, email := range emails {
		details, err := tempmail.GetEmailDetails(email)
		if err != nil {
			t.Errorf("error getting email details %s", err)
		}
		log.Print(details)

		// Get attachments
		for _, attachment := range details.Attachments {
			log.Print(attachment)
		}
	}
}
