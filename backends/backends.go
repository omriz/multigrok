package backends

// Backend defines an interface to a single source fetching backend.
type Backend interface {
	// Query sends a query string to the backend service.
	Query(q string) (WebServiceResult, error)

	// Fetch - fetches a resource from the server. The prefix states what is the fetch mode:
	// xref - annotated html.
	// raw - raw text format.
	Fetch(prefix, path string) ([]byte, error)

	// UID - returns a unique identifier for this backend.
	UID() string
}

// Note - if this ever gets more complicated we'll have to add a deep copy function.
type QueryResult struct {
	Path      string
	Filename  string
	Lineno    int `json:"string,omitempty"`
	Line      string
	Directory string
}

// This is a representation of the OpenGrok Web services response.
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
