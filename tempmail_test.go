package gotempmail

import (
	"fmt"
	"testing"
	"time"
)

func TestTempMailNew(t *testing.T) {
	tempmail := New().
		Address("Testing123@test.com").
		Password("testing123")
	if tempmail == nil {
		t.Error("tempmail is nil")
	}
}

func TestTempMailValidate(t *testing.T) {
	tempmail := New().
		Address("Testing123@test.com").
		Password("testing123")
	if tempmail == nil {
		t.Error("tempmail is nil")
	}

	err := tempmail.Validate()
	if err != nil {
		t.Errorf("validation failed %s", err)
	}
}

func TestCreateAccount(t *testing.T) {
	domains, err := GetDomains()
	if err != nil {
		t.Error(err)
	}

	if domains == nil {
		t.Error("Nil Domains")
	}
	if len(domains) == 0 {
		t.Error("No domains")
	}

	tempmail := New().
		Address("testing" + fmt.Sprintf("%d",
			time.Now().Unix()) + "@" + domains[0]).
		Password("password123").
		CreateAccount()
	if tempmail.Err != nil {
		t.Errorf("tempmail err is %s", tempmail.Err)
	}
}
