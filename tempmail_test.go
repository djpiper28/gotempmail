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
	ret, err := GetDomains()
	if err != nil {
		t.Errorf("error was not expected %s", err)
	}

	if ret == nil {
		t.Errorf("domains is nil")
	}

	if len(ret) == 0 {
		t.Errorf("no domains :(")
	}
}
