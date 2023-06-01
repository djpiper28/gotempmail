package gotempmail

import (
	"testing"
)

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
