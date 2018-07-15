package frontend

import (
	"net/http"
	"net/url"
)

// Const.
var kKeyWords = [...]string{"def", "symbol", "path", "hist", "q", "defs", "refs", "rawquery", "freetext", "full"}

func generateNewQuery(raw string) url.Values {
	q, err := url.ParseQuery(raw)
	if err != nil {
		return q
	}
	for k := range q {
		found := false
		for _, i := range kKeyWords {
			if i == k {
				found = true
				break
			}
		}
		if !found {
			// Not a known key - assuming this is raw text
			q.Del(k)
			q.Add("freetext", k)
		}
	}
	return q
}

func (m *MultiGrokServer) RawSearchHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if len(q) != 1 {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Error in query string"))
		return
	}
	raw, ok := q["rawquery"]
	if !ok || len(raw) != 1 {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Error in query string"))
		return
	}
	// Handle the raw parsing and redirect
	nq := generateNewQuery(raw[0])
	req.URL.RawQuery = nq.Encode()
	m.SearchHandler(w, req)
}
