package middleware

import (
	"strings"
	"testing"

	"github.com/omriz/multigrok/backends"
)

func TestFetchNoBackends(t *testing.T) {
	servers := make(map[string]backends.Backend)
	_, err := Fetch(servers, "/xref/b/c/d.java")
	if !strings.Contains(err.Error(), "Did not find any appropriate backend") {
		t.Errorf("Fetch did not fail with the appropriate error")
	}
}

func TestDecodeFailure(t *testing.T) {
	servers := make(map[string]backends.Backend)
	servers["a"] = &backends.FakeBackend{Id: "a"}
	_, err := Fetch(servers, "/xref/b/c/d.java")
	if !strings.Contains(err.Error(), "Did not find any appropriate backend") {
		t.Errorf("illegal base64 data at input byte")
	}
}
