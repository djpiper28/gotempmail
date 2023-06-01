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
