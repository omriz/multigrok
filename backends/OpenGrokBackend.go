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
	client *http.Client
}

func NewOpenGrokBackend(addr string) OpenGrokBackend {
	opengrokbackend := OpenGrokBackend{
		Addr:   addr,
		client: &http.Client{Timeout: 120 * time.Second},
	}
	return opengrokbackend
}

func (backend *OpenGrokBackend) Query(q string) (QueryResult, error) {
	var result QueryResult
	s := backend.Addr + "json?" + q
	log.Println("Sending request: " + s)
	response, err := backend.client.Get(s)
	if err != nil {
		return result, err
	}
	if response.ContentLength == 0 {
		return result, fmt.Errorf("Malformed request")
	}
	defer response.Body.Close()

	_ = json.NewDecoder(response.Body).Decode(result)
	return result, nil
	// Should now parse the json and get back the results:
	// {"duration":2,"resultcount":0,"freetext":"pytb","results":[]}

}
