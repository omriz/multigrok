package frontend

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

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
	FilePath    string
}

// SearchResultData is a container to display in the rendered results
type searchResultData struct {
	Query        string
	Results      []*fileResult
	TotalResults int
	TimeSecs     float64
}

// Giving priority to code files as this is a code search.
var codeSuffixes = [...]string{".cc", ".h", ".cpp", ".hpp", ".java", ".py", ".go", ".sh", ".pl", ".c", ".m", ".proto"}
var codeFactor = 10

// TODO(omriz): Improve this to get better ranking - for example
// Priority if the query string appears in the path.
func orderResults(resMap map[string]*fileResult) []*fileResult {
	r := make([]*fileResult, 0)
	for len(resMap) > 0 {
		m := 0
		var i string
		for p, r := range resMap {
			w := len(r.LineResults)
			for _, s := range codeSuffixes {
				if strings.HasSuffix(p, s) {
					w += codeFactor
				}
			}
			if w > m {
				m = w
				i = p
			}
		}
		r = append(r, resMap[i])
		delete(resMap, i)
	}
	return r
}

func restructreResults(query string, res backends.WebServiceResult) searchResultData {
	q, err := url.QueryUnescape(query)
	if err != nil {
		q = query
	}
	searchRes := searchResultData{
		Query:        q,
		Results:      make([]*fileResult, 0),
		TotalResults: res.Resultcount,
	}
	resMap := make(map[string]*fileResult)
	for p, first := range res.Results {
		for _, r := range first {
			if val, ok := resMap[p]; ok {
				// We already found this file
				z, err := r.DecodeLine()
				if err == nil {
					resMap[p].LineResults = append(val.LineResults, &lineResult{
						Lineno: r.LineNumber,
						Line:   z,
					})
				}
			} else {
				// New file
				z, err := r.DecodeLine()
				if err == nil {
					pp := strings.Split(p, middleware.SEPERATOR)
					resMap[p] = &fileResult{
						LineResults: []*lineResult{&lineResult{
							Lineno: r.LineNumber,
							Line:   z,
						}},
						ServerPath: p,
						Filename:   filepath.Base(p),
						FilePath:   pp[len(pp)-1],
					}
				}
			}
		}
	}
	searchRes.Results = orderResults(resMap)
	return searchRes
}

func (m *MultiGrokServer) SearchHandler(w http.ResponseWriter, req *http.Request) {
	qparams := req.URL.Query()
	start := time.Now()
	results := middleware.Search(req.Context(), m.backends, qparams)
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
			// The key here is the filepath
			var k []string
			for m := range combined.Results {
				k = append(k, m)
			}
			if combined.Resultcount == 1 {
				p := "xref" + k[0] + "#" + combined.Results[k[0]][0].LineNumber
				log.Println("One result found, redirecting to: " + p)
				http.Redirect(w, req, p, 303)
				return
			}
			// TODO(omriz): Join results under the same file path but different lines to be in the same card.
			data := restructreResults(qparams.Encode(), combined)
			end := time.Now()
			data.TimeSecs = end.Sub(start).Seconds()
			if err := m.resultTmpl.Execute(w, data); err != nil {
				log.Println(err)
			}
		}
	}
}
