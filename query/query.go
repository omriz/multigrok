package main

import (
	"flag"
	"log"
	"strings"

	"github.com/omriz/multigrok/backends"
)

func main() {
	backendAddress := flag.String("backend", "http://localhost:8080/source", "Address to query on")
	query := flag.String("query", "freetext=Hello", "The query text")
	flag.Parse()
	var a string
	if strings.HasSuffix(*backendAddress, "/") {
		a = *backendAddress
	} else {
		a = *backendAddress + "/"
	}
	backend := backends.NewOpenGrokBackend(a)
	log.Println(backend)
	if out, err := backend.Query(*query); err != nil {
		log.Fatalln(err)
	} else {
		log.Println(out)
	}
}
