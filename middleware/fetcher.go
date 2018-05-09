package middleware

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru"

	"golang.org/x/sync/errgroup"

	"github.com/omriz/multigrok/backends"
)

// Fetch finds the appropriate backend(s) to the relevant query and sends the query to them.
// The backend servers are passed as a map UID->Backend
// The query is a full path such as /xref/__MULTIGROKBACKEND__ajlkjhasdlkjh__MULTIGROKBACKEND__/a/b/c.java
func Fetch(servers map[string]backends.Backend, query string, cache *lru.Cache) ([]byte, error) {
	split := strings.Split(query, "/")
	for len(split) > 0 && split[0] == "" {
		split = split[1:]
	}
	if len(split) < 3 {
		if len(split) == 2 {
			return fallBackFetch(servers, split[0], split[1], cache)
		}
		return nil, fmt.Errorf("Can not send request: %v", query)
	}
	cmd := split[0]
	buid, err := DecodeBackendAddress(split[1])
	if err != nil {
		return fallBackFetch(servers, cmd, strings.Join(split[1:], "/"), cache)
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
func fallBackFetch(servers map[string]backends.Backend, cmd, query string, cache *lru.Cache) ([]byte, error) {
	root := strings.Split(query, "/")[0]
	if cache != nil {
		b, ok := cache.Get(root)
		if ok {
			res, err := servers[b.(string)].Fetch(cmd, query)
			if err == nil {
				return res, nil
			}
		}
	}
	return parallelFetch(servers, cmd, query, cache)
}

func parallelFetch(servers map[string]backends.Backend, cmd, query string, cache *lru.Cache) ([]byte, error) {
	type Result struct {
		mu  sync.Mutex
		ok  bool
		res []byte
	}
	res := &Result{}
	g := errgroup.Group{}
	for _, backend := range servers {
		backend := backend // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			ret, err := backend.Fetch(cmd, query)
			if err != nil {
				return nil
			}
			res.mu.Lock()
			defer res.mu.Unlock()
			if !res.ok {
				res.res = ret
				res.ok = true
				// Adding result to the cache.
				if cache != nil {
					root := strings.Split(query, "/")[0]
					cache.Add(root, backend.UID())
				}
			}
			return errors.New("cancel errorgroup")
		})
	}
	g.Wait()
	res.mu.Lock()
	defer res.mu.Unlock()
	if !res.ok {
		return nil, errors.New("Failed to fetch " + query)
	}
	return res.res, nil
}
