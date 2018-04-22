package backends

type BackendType string

type Backend interface {
	Query(q string) (SearchResult, error)
}

type SearchResult struct {
	path      string
	filename  string
	lineno    int64
	line      string
	directory string
}
