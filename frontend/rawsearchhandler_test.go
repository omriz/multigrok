package frontend

import "testing"

func TestQueryDecode(t *testing.T) {
	q := "abcd&defs=blabla"
	nq := generateNewQuery(q)
	if len(nq) != 2 {
		t.Errorf("Too many query params %v", nq)
	}
	if len(nq["freetext"]) != 1 || nq["freetext"][0] != "abcd" {
		t.Errorf("Missing freetext 'abcd' in %v", nq)
	}
	if len(nq["defs"]) != 1 || nq["defs"][0] != "blabla" {
		t.Errorf("Missing defs 'blabla' in %v", nq)
	}
}
