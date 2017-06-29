package user

import (
	"testing"
)

func TestPassword(t *testing.T) {
	p := NewPassword("asdf")

	if !p.Check("asdf") {
		t.Errorf("password not matched was: %s", p)
	}
	if p.Check("fdsa") {
		t.Errorf("password matched all string: %s", p.Data)
	}
}
