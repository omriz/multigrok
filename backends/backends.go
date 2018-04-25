package backends

type BackendType string

type Backend interface {
	Query(q string) (WebServiceResult, error)
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
