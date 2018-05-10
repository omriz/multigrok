package frontend

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
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
	q, err := url.QueryUnescape(query)
	if err != nil {
		q = query
	}
	searchRes := searchResultData{
		Query:        q,
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
	results := middleware.Search(m.backends, qparams)
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
				p := "xref" + combined.Results[0].Path + "#" + combined.Results[0].Lineno
				log.Println("One result found, redirecting to: " + p)
				http.Redirect(w, req, p, 303)
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
