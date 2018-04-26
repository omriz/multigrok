package main

import (
	"flag"
	"log"
	"strings"

	"github.com/omriz/multigrok/backends"
	"github.com/omriz/multigrok/middleware"
)

func main() {
	backendsFlag := flag.String("backends", "http://localhost:8080/source", "Comma seperated list of backends")
	fpath := flag.String("fpath", "/xref/bla.c", "The query text")
	flag.Parse()
	ogbs := make(map[string]backends.Backend)
	for _, a := range strings.Split(*backendsFlag, ",") {
		if !strings.HasSuffix(a, "/") {
			a = a + "/"
		}
		b := backends.NewOpenGrokBackend(a)
		ogbs[b.UID()] = &b
	}
	data, err := middleware.Fetch(ogbs, *fpath)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	log.Println(string(data))
}
