package backends

type BackendType string

type Backend interface {
	Query(q string) (QueryResult, error)
}

// This is a representation of the OpenGrok Web services response.
// An example can be seen at:
// https://github.com/oracle/opengrok/wiki/Web-services
type QueryResult struct {
	Duration    int64
	Path        string
	Resultcount int
	Hist        string
	Freetext    string
	Results     []struct {
		Path      string
		Filename  string
		Lineno    int
		Line      string
		Directory string
	}
}
