package frontend

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/omriz/multigrok/backends"
)

type MultiGrokServer struct {
	backends         map[string]backends.Backend
	port             int
	client           *http.Client
	loopbackPrefixes []string
	resultTmpl       *template.Template
}

func NewMultiGrokServer(backends map[string]backends.Backend, port int) *MultiGrokServer {
	mgs := MultiGrokServer{
		backends:         backends,
		port:             port,
		client:           &http.Client{Timeout: 120 * time.Second},
		loopbackPrefixes: []string{"/source"},
		resultTmpl:       template.Must(template.ParseFiles("frontend/templates/results.html")),
	}
	http.Handle("/default/", http.FileServer(http.Dir("static")))
	http.Handle("/js/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/source/", mgs.LoopBackHandler)
	http.HandleFunc("/rawsearch", mgs.RawSearchHandler)
	http.HandleFunc("/search", mgs.SearchHandler)
	http.HandleFunc("/s", mgs.SearchHandler)
	http.HandleFunc("/xref/", mgs.FetchHandler)
	http.HandleFunc("/raw/", mgs.FetchHandler)
	http.HandleFunc("/history/", mgs.FetchHandler)
	http.HandleFunc("/diff/", mgs.FetchHandler)
	http.HandleFunc("/download/", mgs.FetchHandler)
	http.HandleFunc("/index.html", mgs.RootHandler)
	http.HandleFunc("/", mgs.RootHandler)
	return &mgs
}

func (m *MultiGrokServer) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
}
