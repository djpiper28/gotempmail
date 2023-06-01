package gotempmail

import (
	"testing"
)

func TestTempMailNew(t *testing.T) {
	tempmail := New()
	if tempmail == nil {
		t.Error("tempmail is nil")
	}
}

func TestGetDomainsHasNoErr(t *testing.T) {
	_, err := GetDomains()
	if err != nil {
		t.Errorf("error was not expected %e", err)
	}
}
