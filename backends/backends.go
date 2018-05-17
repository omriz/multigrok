package backends

import (
	"context"
	"encoding/base64"
	"html/template"
	"strings"
)

// Backend defines an interface to a single source fetching backend.
type Backend interface {
	// Query sends a query string to the backend service.
	Query(ctx context.Context, q string) (WebServiceResult, error)

	// Fetch - fetches a resource from the server. The prefix states what is the fetch mode:
	// xref - annotated html.
	// raw - raw text format.
	Fetch(ctx context.Context, prefix, path string) ([]byte, error)

	// UID - returns a unique identifier for this backend.
	UID() string
}

// QueryResult represents a single search result..
type QueryResult struct {
	Path      string
	Filename  string
	Lineno    string
	Line      string
	Directory string
}

// DecodeLine used to decode the base64 encoded html string.
// Returns a template.HTML in order to avoid double escaping the response.
func (q *QueryResult) DecodeLine() (template.HTML, error) {
	d, err := base64.StdEncoding.DecodeString(q.Line)
	if err != nil {
		return "", err
	}
	return template.HTML(d), nil
}

// RefactorPath removes the escaping from the directory structure.
func (q *QueryResult) RefactorPath() string {
	return strings.Replace(q.Directory, "\\/", "/", -1)
}

// WebServiceResult is a representation of the OpenGrok Web services response.
// An example can be seen at:
// https://github.com/oracle/opengrok/wiki/Web-services
type WebServiceResult struct {
	Duration    int64
	Path        string
	Resultcount int
	Hist        string
	Freetext    string
	Results     []QueryResult
}
