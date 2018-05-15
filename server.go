package main

import (
	"flag"
	"io/ioutil"
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
	httpPort := flag.Int("http_port", 80, "port to run the http web service on if avaliable")
	httpsPort := flag.Int("https_port", 443, "port to run the https web service on if avaliebl")
	certFile := flag.String("cert_file", "", "If using https, the path to the certificate file")
	keyFile := flag.String("key_file", "", "If using https, the path to the key file")
	hostName := flag.String("hostname", "", "Your hostname if using auto cert renewal")
	mode := flag.String("mode", "http", "One of the following: http,https, autoCert")

	flag.Parse()
	// Validating flags
	if !(*mode == "http" || *mode == "https" || *mode != "autoCert") {
		log.Fatalf("mode %s is invalid. Avaliable modes are http, https and autoCert", *mode)
	}
	ogbs := make(map[string]backends.Backend)
	for _, a := range strings.Split(*backendsFlag, ",") {
		if !strings.HasSuffix(a, "/") {
			a = a + "/"
		}
		ogb := backends.NewOpenGrokBackend(a)
		ogbs[ogb.UID()] = &ogb
	}
	s := frontend.NewMultiGrokServer(ogbs, *httpPort, *httpsPort)
	if *mode == "https" {
		if *certFile == "" || *keyFile == "" {
			log.Fatalf("Missing certification or key files")
		}
		log.Fatal(s.ListenAndServeHttps(*certFile, *keyFile))
	} else if *mode == "autoCert" {
		if *hostName == "" {
			log.Fatalf("Missing host_name")
		}
		cacheDir, err := ioutil.TempDir("", "certs")
		if err != nil {
			log.Fatalf("Failed to created certificate cache dir")
		}
		log.Fatal(s.ListenAndServeAutoCert(*hostName, cacheDir))
	} else {
		log.Fatal(s.ListenAndServeHttp())
	}
}
