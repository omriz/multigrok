package backends

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// OpenGrokBackend is a representation of an OpenGrok server backend.
type OpenGrokBackend struct {
	addr   string
	client *http.Client
}

// NewOpenGrokBackend creates a new OpenGrok backend.
func NewOpenGrokBackend(addr string) OpenGrokBackend {
	var a string
	if !strings.HasSuffix(addr, "/") {
		a = addr + "/"
	} else {
		a = addr
	}
	opengrokbackend := OpenGrokBackend{
		addr:   a,
		client: &http.Client{Timeout: 120 * time.Second},
	}
	return opengrokbackend
}

// Query sends a query to our backend
func (backend *OpenGrokBackend) Query(ctx context.Context, q string) (WebServiceResult, error) {
	var result WebServiceResult
	s := backend.addr + "api/v1/search?" + q
	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		log.Printf("Failed to initialize a new request %s", s)
		return result, err
	}
	req = req.WithContext(ctx)
	response, err := backend.client.Do(req)
	if err != nil {
		log.Printf("Got error: %v\n", err)
		return result, err
	}
	if response.ContentLength == 0 {
		return result, fmt.Errorf("Malformed request")
	}
	defer response.Body.Close()
	temp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(temp, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UID returns the backend address as its identifier
func (backend *OpenGrokBackend) UID() string {
	return backend.addr
}

// Fetch - returns a resource
func (backend *OpenGrokBackend) Fetch(ctx context.Context, prefix, path string) ([]byte, error) {
	a := strings.TrimPrefix(path, "/")
	s := backend.addr + prefix + "/" + a
	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		log.Printf("Failed to initialize a new request %s", s)
		return nil, err
	}
	req = req.WithContext(ctx)
	response, err := backend.client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.ContentLength == 0 {
		return nil, fmt.Errorf("Malformed request")
	}
	if response.StatusCode > 400 {
		return nil, fmt.Errorf("Error fetching page %v", response.Status)
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	return result, nil
}
