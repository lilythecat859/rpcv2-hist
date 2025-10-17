package rest

import (
	"net/http"

	"github.com/faithful-rpc/rpcv2-hist/internal/health"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health.NewHandler(s.root, s.log).ServeHTTP(w, r)
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	// same as health for now
	s.handleHealth(w, r)
}