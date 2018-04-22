package backends

import (
	"fmt"
	"log"
	"net/http"
)

type BackendType string

const (
	OPENGROK BackendType = "OPENGROK"
)

type Backend struct {
	Addr string
	Type BackendType
}

type SearchResult struct {
	path      string
	filename  string
	lineno    int64
	line      string
	directory string
}

func (backend *Backend) Query(q string) (SearchResult, error) {
	s := backend.Addr + "json?" + q
	log.Println("Sending request: " + s)
	res, err := http.Get(s)
	if res.ContentLength == 0 {
		return nil, fmt.Errorf("Malformed request")
	}
	// Should now parse the json and get back the results:
	// {"duration":2,"resultcount":0,"freetext":"pytb","results":[]}

}
