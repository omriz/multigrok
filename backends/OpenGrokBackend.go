package backends

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OpenGrokBackend struct {
	Addr   string
	client *http.Client
}

func NewOpenGrokBackend(addr string) OpenGrokBackend {
	var a string
	if !strings.HasSuffix(addr, "/") {
		a = addr + "/"
	} else {
		a = addr
	}
	opengrokbackend := OpenGrokBackend{
		Addr:   a,
		client: &http.Client{Timeout: 120 * time.Second},
	}
	return opengrokbackend
}

func (backend *OpenGrokBackend) Query(q string) (QueryResult, error) {
	var result QueryResult
	s := backend.Addr + "json?" + url.QueryEscape(q)
	log.Println("Sending request: " + s)
	response, err := backend.client.Get(s)
	if err != nil {
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
