package backends

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestA(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		qr := QueryResult{
			Duration:    1000,
			Path:        "",
			Resultcount: 0,
			Freetext:    "TestTest",
		}
		json.NewEncoder(w).Encode(qr)
	}))
	defer ts.Close()
	b := NewOpenGrokBackend(ts.URL)
	res, err := b.Query("freetext=Ho")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
