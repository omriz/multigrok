package middleware

import (
	"testing"

	"github.com/omriz/multigrok/backends"
)

func TestEncodeBackendAddress(t *testing.T) {
	orig := "http://test.multigrok.backend/"
	encoded := EncodeBackendAddress(orig)
	expected := SEPERATOR + "aHR0cDovL3Rlc3QubXVsdGlncm9rLmJhY2tlbmQv" + SEPERATOR
	if encoded != expected {
		t.Errorf("Expected: " + expected)
		t.Errorf("Result: " + encoded)
	}
}

func TestDecodeBackendAddress(t *testing.T) {
	expected := "http://test.multigrok.backend/"
	// Intenitonally leaving out the SEPERATOR to verify it doesn't change.
	encoded := "__MULTIGROKBACKEND__aHR0cDovL3Rlc3QubXVsdGlncm9rLmJhY2tlbmQv__MULTIGROKBACKEND__"
	decoded, err := DecodeBackendAddress(encoded)
	if err != nil {
		t.Fatalf("Got Error: %v", err)
	}
	if decoded != expected {
		t.Errorf("Expected: " + expected)
		t.Errorf("Result: " + decoded)
	}
}

func compareSlices(a, b []backends.QueryResult) bool {
	if len(a) != len(b) {
		return false
	}
	for _, elem := range a {
		found := false
		for _, x := range b {
			if elem == x {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func TestCombineResults(t *testing.T) {
	responses := make(map[string]backends.WebServiceResult)
	responses["http://1.1.1.1/"] = backends.WebServiceResult{
		Resultcount: 2,
		Results: []backends.QueryResult{
			{
				Path:      "/a/b/c.java",
				Filename:  "c.java",
				Directory: "\\/a\\/b\\/",
				Lineno:    "8",
			},
			{
				Path:      "/d/e.java",
				Filename:  "e.java",
				Directory: "\\/d\\/",
				Lineno:    "1",
			},
		},
	}
	responses["http://1.1.1.2/"] = backends.WebServiceResult{
		Resultcount: 1,
		Results: []backends.QueryResult{
			{
				Path:      "/xx/yy/c.java",
				Filename:  "c.java",
				Directory: "\\/xx\\/yy\\/",
				Lineno:    "8",
			},
		},
	}
	responses["http://1.1.1.3/"] = backends.WebServiceResult{
		Resultcount: 0,
	}
	expected := backends.WebServiceResult{
		Resultcount: 3,
		Results: []backends.QueryResult{
			{
				Path:      "/" + EncodeBackendAddress("http://1.1.1.1/") + "/a/b/c.java",
				Filename:  "c.java",
				Directory: "\\/a\\/b\\/",
				Lineno:    "8",
			},
			{
				Path:      "/" + EncodeBackendAddress("http://1.1.1.1/") + "/d/e.java",
				Filename:  "e.java",
				Directory: "\\/d\\/",
				Lineno:    "1",
			},
			{
				Path:      "/" + EncodeBackendAddress("http://1.1.1.2/") + "/xx/yy/c.java",
				Filename:  "c.java",
				Directory: "\\/xx\\/yy\\/",
				Lineno:    "8",
			},
		},
	}
	combined, err := CombineResults(responses)
	if err != nil {
		t.Fatalf("Got Error: %v", err)
	}
	if expected.Resultcount != combined.Resultcount || !compareSlices(expected.Results, combined.Results) {
		t.Errorf("Expected: %v", expected)
		t.Errorf("Result: %v", combined)
	}
}
