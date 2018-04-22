package backends

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type OpenGrokBackend struct {
	Addr   string
	client http.Client
}

func NewOpenGrokBackend(addr string) OpenGrokBackend {
	opengrokbackend := OpenGrokBackend{
		Addr:   addr,
		client: &http.Client{Timeout: 120 * time.Second},
	}
	return opengrokbackend
}

func (backend *OpenGrokBackend) Query(q string) (SearchResult, error) {
	s := backend.Addr + "json?" + q
	log.Println("Sending request: " + s)
	res, err := client.Get(s)
	if err != nil {
		return nil, err
	}
	if res.ContentLength == 0 {
		return nil, fmt.Errorf("Malformed request")
	}
	defer r.Body.Close()

	j := json.NewDecoder(r.Body).Decode(target)
	// Should now parse the json and get back the results:
	// {"duration":2,"resultcount":0,"freetext":"pytb","results":[]}

}
