package internal

import (
	"testing"
)

func AssertEqualString(t *testing.T, name, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("%s: got %q, expected %q", name, got, want)
	}
}
