package frontend

import (
	"net/http"

	"github.com/omriz/multigrok/backends"
)

type MultiGrokServer struct {
	backends map[string]backends.Backend
}

func NewMultiGrokServer(backends map[string]backends.Backend) *MultiGrokServer {
	mgs := MultiGrokServer{
		backends: backends,
	}
	http.HandleFunc("/search", mgs.SearchHandler)
	return &mgs
}

func (m *MultiGrokServer) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
