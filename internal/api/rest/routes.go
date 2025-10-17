package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/block/{slot}", s.handleGetBlock).Methods("GET")
	r.HandleFunc("/tx/{signature}", s.handleGetTx).Methods("GET")
	r.HandleFunc("/sigs/{address}", s.handleGetSigs).Methods("GET")
	r.HandleFunc("/health", s.handleHealth).Methods("GET")
	r.HandleFunc("/ready", s.handleReady).Methods("GET")
	return r
}