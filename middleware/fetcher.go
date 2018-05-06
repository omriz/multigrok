package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/omriz/multigrok/backends"
)

// Fetch finds the appropriate backend(s) to the relevant query and sends the query to them.
// The backend servers are passed as a map UID->Backend
// The query is a full path such as /xref/__MULTIGROKBACKEND__ajlkjhasdlkjh__MULTIGROKBACKEND__/a/b/c.java
func Fetch(servers map[string]backends.Backend, query string) ([]byte, error) {
	split := strings.Split(query, "/")
	for len(split) > 0 && split[0] == "" {
		split = split[1:]
	}
	if len(split) < 3 {
		if len(split) == 2 {
			return fallBackFetch(servers, split[0], split[1])
		}
		return nil, fmt.Errorf("Can not send request: %v", query)
	}
	cmd := split[0]
	log.Printf("Cmd: '%s'\n", cmd)
	log.Printf("Encoded Backend: '%s'\n", split[1])
	buid, err := DecodeBackendAddress(split[1])
	if err != nil {
		return fallBackFetch(servers, cmd, strings.Join(split[1:], "/"))
	}
	for uid, backend := range servers {
		if buid == uid {
			return backend.Fetch(cmd, strings.Join(split[2:], "/"))
		}
	}
	// TODO(omriz): We should have a code here that attemps direct fetching for a
	// path with no backends.
	return nil, fmt.Errorf("Did not find any appropriate backend")
}

// Attempts to fetch from each backend - should be done in parallel?
func fallBackFetch(servers map[string]backends.Backend, cmd, query string) ([]byte, error) {
	for _, backend := range servers {
		x, err := backend.Fetch(cmd, query)
		if err == nil {
			return x, nil
		}
	}
	return nil, fmt.Errorf("Could not fetch" + query)
}
