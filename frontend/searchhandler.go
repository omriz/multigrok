package frontend

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/omriz/multigrok/backends"
	"github.com/omriz/multigrok/middleware"
)

func (m *MultiGrokServer) SearchHandler(w http.ResponseWriter, req *http.Request) {
	qparams := req.URL.Query()
	if val, ok := qparams["q"]; ok {
		qparams.Add("freetext", val[0])
		qparams.Del("q")
	}
	// TODO(omriz): This should be made parallel in the future.
	results := make(map[string]backends.WebServiceResult)
	for name, backend := range m.backends {
		res, err := backend.Query(qparams.Encode())
		if err != nil {
			log.Printf("Error returned from %s: %v\n", name, err)
		} else {
			// We only append the result if we have any.
			results[name] = res
		}
	}
	if len(results) == 0 {
		log.Printf("No results found")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("No results found for query: %v.\n", qparams)))
	} else {
		combined, err := middleware.CombineResults(results)
		if err != nil {
			log.Printf("Failed combining")
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(fmt.Sprintf("Error parsing results for query: %v.\n%v", qparams, err)))
		} else {
			w.Header().Set("Content-Type", "text/html")
			c := fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head>\n<title>Results: %s</title>\n</head>\n<body><h1>Response</h1>", qparams)
			for _, res := range combined.Results {
				content, err := base64.StdEncoding.DecodeString(res.Line)
				if err != nil {
					c += fmt.Sprintf("<b>%s></b><br>%s", res.Path, "error")
				} else {
					c += fmt.Sprintf("<b>%s></b><br>%s", res.Path, content)
				}
			}
			c += "</body></html>"
			log.Printf("responding with %v", c)
			w.Write([]byte(c))
		}
	}
}
