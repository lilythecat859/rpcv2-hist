package jsonrpc

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	r := mux.NewRouter()
	r.Handle("/", s).Methods("POST")
	r.HandleFunc("/health", s.handleHealth).Methods("GET")
	return r
}