package frontend

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/omriz/multigrok/backends"
	"github.com/omriz/multigrok/middleware"
)

type lineResult struct {
	Lineno string
	Line   template.HTML
}
type fileResult struct {
	ServerPath  string
	Filename    string
	LineResults []*lineResult
}

// SearchResultData is a container to display in the rendered results
type searchResultData struct {
	Query        string
	Results      map[string]*fileResult
	TotalResults int
}

func restructreResults(query string, res backends.WebServiceResult) searchResultData {
	searchRes := searchResultData{
		Query:        query,
		Results:      make(map[string]*fileResult),
		TotalResults: res.Resultcount,
	}
	for _, first := range res.Results {
		p := filepath.Join(first.RefactorPath(), first.Filename)
		if val, ok := searchRes.Results[p]; ok {
			// We already found this file
			z, err := first.DecodeLine()
			if err == nil {
				searchRes.Results[p].LineResults = append(val.LineResults, &lineResult{
					Lineno: first.Lineno,
					Line:   z,
				})
			}
		} else {
			// New file
			z, err := first.DecodeLine()
			if err == nil {
				searchRes.Results[p] = &fileResult{
					LineResults: []*lineResult{&lineResult{
						Lineno: first.Lineno,
						Line:   z,
					}},
					ServerPath: first.Path,
					Filename:   first.Filename,
				}
			}
		}
	}
	return searchRes
}

func (m *MultiGrokServer) SearchHandler(w http.ResponseWriter, req *http.Request) {
	qparams := req.URL.Query()
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
			if combined.Resultcount == 1 {
				log.Println("One result found, redirecting to: xref" + combined.Results[0].Path)
				http.Redirect(w, req, "xref"+combined.Results[0].Path, 303)
				return
			}
			// TODO(omriz): Join results under the same file path but different lines to be in the same card.
			data := restructreResults(qparams.Encode(), combined)
			if err := m.resultTmpl.Execute(w, data); err != nil {
				log.Println(err)
			}
		}
	}
}
