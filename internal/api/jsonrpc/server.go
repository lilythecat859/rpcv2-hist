package jsonrpc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/lilythecat859/rpcv2-hist/internal/fractal"
	"github.com/lilythecat859/rpcv2-hist/internal/model"
	"github.com/lilythecat859/rpcv2-hist/internal/storage"
)

const (
	// Solana JSON-RPC spec
	version = "2.0"
)

type Server struct {
	root   *fractal.Root
	log    *zap.Logger
	tracer trace.Tracer
}

type request struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type response struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	errInvalidRequest = &rpcError{-32600, "Invalid request"}
	errMethodNotFound = &rpcError{-32601, "Method not found"}
	errInternal       = &rpcError{-32603, "Internal error"}
)

func NewServer(root *fractal.Root, log *zap.Logger) http.Handler {
	s := &Server{
		root:   root,
		log:    log,
		tracer: otel.Tracer("jsonrpc"),
	}
	r := mux.NewRouter()
	r.Handle("/", s).Methods("POST")
	return r
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx, span := s.tracer.Start(r.Context(), "JSON-RPC", trace.WithAttributes(attribute.String("method", r.Method)))
	defer span.End()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, errInvalidRequest, nil)
		return
	}
	span.SetAttributes(attribute.String("rpc.method", req.Method))

	var result interface{}
	var rpcErr *rpcError
	switch strings.ToLower(req.Method) {
	case "getblock":
		result, rpcErr = s.handleGetBlock(ctx, req.Params)
	case "gettransaction":
		result, rpcErr = s.handleGetTransaction(ctx, req.Params)
	case "getsignaturesforaddress":
		result, rpcErr = s.handleGetSignaturesForAddress(ctx, req.Params)
	case "getblockswithlimit":
		result, rpcErr = s.handleGetBlocksWithLimit(ctx, req.Params)
	case "getblocktime":
		result, rpcErr = s.handleGetBlockTime(ctx, req.Params)
	default:
		rpcErr = errMethodNotFound
	}

resp := response{
		Jsonrpc: version,
		ID:      req.ID,
	}
	if rpcErr != nil {
		resp.Error = rpcErr
	} else {
		resp.Result = result
	}

	s.writeJSON(w, resp)

	s.log.Info("request",
		zap.String("method", req.Method),
		zap.Duration("dur", time.Since(start)),
		zap.Bool("error", rpcErr != nil),
	)
}

func (s *Server) handleGetBlock(ctx context.Context, params json.RawMessage) (interface{}, *rpcError) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil || len(p) < 1 {
		return nil, errInvalidRequest
	}
	slot, ok := toUint64(p[0])
	if !ok {
		return nil, errInvalidRequest
	}
	commit := commitmentFromParams(p, 1)
	blk, err := s.root.GetBlock(ctx, slot, commit)
	if err != nil {
		return nil, errInternal
	}
	return blk, nil
}

func (s *Server) handleGetTransaction(ctx context.Context, params json.RawMessage) (interface{}, *rpcError) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil || len(p) < 1 {
		return nil, errInvalidRequest
	}
	sig, ok := p[0].(string)
	if !ok {
		return nil, errInvalidRequest
	}
	commit := commitmentFromParams(p, 1)
	tx, err := s.root.GetTransaction(ctx, sig, commit)
	if err != nil {
		return nil, errInternal
	}
	return tx, nil
}

func (s *Server) handleGetSignaturesForAddress(ctx context.Context, params json.RawMessage) (interface{}, *rpcError) {
	var p map[string]interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, errInvalidRequest
	}
	addr, _ := p["address"].(string)
	if addr == "" {
		return nil, errInvalidRequest
	}
	opts := storage.SignatureOpts{
		Limit:      1000,
		Commitment: storage.CommitmentConfirmed,
	}
	if l, ok := toUint64(p["limit"]); ok && l > 0 && l <= 1000 {
		opts.Limit = l
	}
	if b, ok := p["before"].(string); ok {
		opts.Before = &b
	}
	if u, ok := p["until"].(string); ok {
		opts.Until = &u
	}
	sigs, err := s.root.GetSignaturesForAddress(ctx, addr, opts)
	if err != nil {
		return nil, errInternal
	}
	return sigs, nil
}

func (s *Server) handleGetBlocksWithLimit(ctx context.Context, params json.RawMessage) (interface{}, *rpcError) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil || len(p) < 2 {
		return nil, errInvalidRequest
	}
	start, ok1 := toUint64(p[0])
	limit, ok2 := toUint64(p[1])
	if !ok1 || !ok2 || limit == 0 || limit > 500000 {
		return nil, errInvalidRequest
	}
	commit := commitmentFromParams(p, 2)
	slots, err := s.root.GetBlocksWithLimit(ctx, start, limit, commit)
	if err != nil {
		return nil, errInternal
	}
	return slots, nil
}

func (s *Server) handleGetBlockTime(ctx context.Context, params json.RawMessage) (interface{}, *rpcError) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil || len(p) < 1 {
		return nil, errInvalidRequest
	}
	slot, ok := toUint64(p[0])
	if !ok {
		return nil, errInvalidRequest
	}
	t, err := s.root.GetBlockTime(ctx, slot)
	if err != nil {
		return nil, errInternal
	}
	if t == nil {
		return nil, nil
	}
	return t.Unix(), nil
}

func (s *Server) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) writeError(w http.ResponseWriter, err *rpcError, id interface{}) {
	resp := response{Jsonrpc: version, ID: id, Error: err}
	s.writeJSON(w, resp)
}

func toUint64(v interface{}) (uint64, bool) {
	switch n := v.(type) {
	case float64:
		return uint64(n), true
	case int:
		return uint64(n), true
	case uint64:
		return n, true
	default:
		return 0, false
	}
}

func commitmentFromParams(p []interface{}, idx int) storage.Commitment {
	if idx < len(p) {
		if m, ok := p[idx].(map[string]interface{}); ok {
			if c, ok := m["commitment"].(string); ok {
				return storage.Commitment(c)
			}
		}
	}
	return storage.CommitmentFinalized
}