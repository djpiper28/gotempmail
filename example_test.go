package gotempmail

import (
	"fmt"
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
}
