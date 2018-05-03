package frontend

import "net/http"

func (m *MultiGrokServer) RootHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "static/index.html")
}
