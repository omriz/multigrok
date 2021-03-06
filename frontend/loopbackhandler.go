package frontend

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (m *MultiGrokServer) LoopBackHandler(w http.ResponseWriter, req *http.Request) {
	//c := strings.TrimPrefix(req.RequestURI, m.port)
	c := req.RequestURI
	for _, x := range m.loopbackPrefixes {
		c = strings.TrimPrefix(c, x)
	}
	nreq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d%s", m.httpPort, c), nil)
	nreq = nreq.WithContext(req.Context())
	resp, err := m.client.Do(nreq)
	if err != nil {
		log.Printf("Failed fetching responses %v\n", err)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Failed fetching responses")))
	} else {
		defer resp.Body.Close()
		temp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(fmt.Sprintf("Failed fetching responses")))
		} else {
			for k, v := range resp.Header {
				for _, v1 := range v {
					w.Header().Add(k, v1)
				}
			}
			w.Write(temp)
		}
	}
}
