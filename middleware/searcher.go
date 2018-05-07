package middleware

import (
	"log"
	"net/url"

	"github.com/omriz/multigrok/backends"
)

func Search(servers map[string]backends.Backend, qparams url.Values) map[string]backends.WebServiceResult {
	results := make(map[string]backends.WebServiceResult)
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
	for name, backend := range servers {
		res, err := backend.Query(qparams.Encode())
		if err != nil {
			log.Printf("Error returned from %s: %v\n", name, err)
		} else {
			// We only append the result if we have any.
			results[name] = res
		}
	}
	return results
}
