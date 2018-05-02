package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/omriz/multigrok/backends"
	"github.com/omriz/multigrok/frontend"
)

func searchPlaceholder(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is a place holder for search server.\n"))
}

func main() {
	// Flags definitions.
	backendsFlag := flag.String("backends", "http://localhost:8080/source", "Comma seperated list of backends")
	port := flag.Int("port", 8080, "port to run the web service on")
	flag.Parse()
	ogbs := make(map[string]backends.Backend)
	for _, a := range strings.Split(*backendsFlag, ",") {
		if !strings.HasSuffix(a, "/") {
			a = a + "/"
		}
		ogb := backends.NewOpenGrokBackend(a)
		ogbs[ogb.UID()] = &ogb
	}
	s := frontend.NewMultiGrokServer(ogbs, *port)
	// Registring and running the server.
	log.Fatal(s.ListenAndServe())
}
