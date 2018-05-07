package frontend

import (
	"fmt"
	"net/http"

	"github.com/omriz/multigrok/middleware"
)

func (m *MultiGrokServer) FetchHandler(w http.ResponseWriter, req *http.Request) {
	res, err := middleware.Fetch(m.backends, req.URL.RequestURI())
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Failed fetching responses")))
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Write(res)
	}
}
