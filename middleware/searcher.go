package middleware

import (
	"log"
	"net/url"
	"sync"

	"github.com/omriz/multigrok/backends"
)

func parseParams(qparams url.Values) string {
	if val, ok := qparams["q"]; ok {
		qparams.Add("freetext", val[0])
		qparams.Del("q")
	}
	if val, ok := qparams["defs"]; ok {
		qparams.Add("def", val[0])
		qparams.Del("defs")
	}
	if val, ok := qparams["refs"]; ok {
		qparams.Add("symbol", val[0])
		qparams.Del("refs")
	}
	if _, ok := qparams["project"]; ok {
		// We intentionally drop the project.
		qparams.Del("project")
	}
	return qparams.Encode()
}

func Search(servers map[string]backends.Backend, qparams url.Values) map[string]backends.WebServiceResult {
	q := parseParams(qparams)
	// Parallel execution over all backends.
	var wg sync.WaitGroup
	var m sync.Map
	for name, backend := range servers {
		wg.Add(1)
		go func(name string, backend backends.Backend, query string, m *sync.Map, wg *sync.WaitGroup) {
			res, err := backend.Query(query)
			if err != nil {
				log.Printf("Error returned from %s: %v\n", name, err)
			} else {
				// We only append the result if we have any.
				m.Store(name, res)
			}
			wg.Done()
		}(name, backend, q, &m, &wg)
	}
	wg.Wait()
	results := make(map[string]backends.WebServiceResult)
	for name := range servers {
		if v, ok := m.Load(name); ok {
			if c, k := v.(backends.WebServiceResult); k {
				results[name] = c
			}
		}
	}
	return results
}
