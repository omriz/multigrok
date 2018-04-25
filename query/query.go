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
	query := flag.String("query", "freetext=Hello", "The query text")
	flag.Parse()
	var ogbs []backends.OpenGrokBackend
	for _, a := range strings.Split(*backendsFlag, ",") {
		if !strings.HasSuffix(a, "/") {
			a = a + "/"
		}
		ogbs = append(ogbs, backends.NewOpenGrokBackend(a))
	}
	results := make(map[string]backends.WebServiceResult)
	for _, ogb := range ogbs {
		if out, err := ogb.Query(*query); err != nil {
			log.Fatalln(err)
		} else {
			log.Printf("Found %d results", out.Resultcount)
			results[ogb.Addr] = out
		}
	}
	combined, err := middleware.CombineResults(results)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Total %d results", combined.Resultcount)
}
