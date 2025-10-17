package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/faithful-rpc/rpcv2-hist/internal/fractal"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage"
)

type Server struct {
	root   *fractal.Root
	log    *zap.Logger
	tracer trace.Tracer
}

func NewServer(root *fractal.Root, log *zap.Logger) http.Handler {
	s := &Server{
		root:   root,
		log:    log,
		tracer: otel.Tracer("rest"),
	}
	r := mux.NewRouter()
	r.HandleFunc("/block/{slot}", s.handleGetBlock).Methods("GET")
	r.HandleFunc("/tx/{signature}", s.handleGetTx).Methods("GET")
	r.HandleFunc("/sigs/{address}", s.handleGetSigs).Methods("GET")
	return r
}

func (s *Server) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	ctx, span := s.tracer.Start(r.Context(), "REST GetBlock")
	defer span.End()

	vars := mux.Vars(r)
	slot, err := strconv.ParseUint(vars["slot"], 10, 64)
	if err != nil {
		http.Error(w, "invalid slot", http.StatusBadRequest)
		return
	}
	commit := storage.Commitment(r.URL.Query().Get("commitment"))
	if commit == "" {
		commit = storage.CommitmentFinalized
	}
	blk, err := s.root.GetBlock(ctx, slot, commit)
	if err != nil {
		s.log.Warn("getBlock", zap.Uint64("slot", slot), zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(blk)
}

func (s *Server) handleGetTx(w http.ResponseWriter, r *http.Request) {
	ctx, span := s.tracer.Start(r.Context(), "REST GetTx")
	defer span.End()

	sig := mux.Vars(r)["signature"]
	commit := storage.Commitment(r.URL.Query().Get("commitment"))
	if commit == "" {
		commit = storage.CommitmentFinalized
	}
	tx, err := s.root.GetTransaction(ctx, sig, commit)
	if err != nil {
		s.log.Warn("getTx", zap.String("sig", sig), zap.Error(err))
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tx)
}

func (s *Server) handleGetSigs(w http.ResponseWriter, r *http.Request) {
	ctx, span := s.tracer.Start(r.Context(), "REST GetSigs")
	defer span.End()

	addr := mux.Vars(r)["address"]
	limit, _ := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if limit == 0 || limit > 1000 {
		limit = 100
	}
	opts := storage.SignatureOpts{
		Limit:      limit,
		Commitment: storage.Commitment(r.URL.Query().Get("commitment")),
	}
	if opts.Commitment == "" {
		opts.Commitment = storage.CommitmentConfirmed
	}
	sigs, err := s.root.GetSignaturesForAddress(ctx, addr, opts)
	if err != nil {
		s.log.Warn("getSigs", zap.String("addr", addr), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sigs)
}