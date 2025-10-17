package grpc

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/faithful-rpc/rpcv2-hist/internal/fractal"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage"
)

type Server struct {
	UnimplementedHistoricalServer
	root   *fractal.Root
	log    *zap.Logger
	tracer trace.Tracer
}

func NewServer(root *fractal.Root, log *zap.Logger) *Server {
	return &Server{
		root:   root,
		log:    log,
		tracer: otel.Tracer("grpc"),
	}
}

func (s *Server) GetBlock(ctx context.Context, req *GetBlockRequest) (*GetBlockResponse, error) {
	blk, err := s.root.GetBlock(ctx, req.Slot, storage.Commitment(req.Commitment))
	if err != nil {
		return nil, err
	}
	return &GetBlockResponse{Raw: blk.Raw}, nil
}

func (s *Server) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	tx, err := s.root.GetTransaction(ctx, req.Signature, storage.Commitment(req.Commitment))
	if err != nil {
		return nil, err
	}
	return &GetTransactionResponse{Raw: tx.Raw}, nil
}

func (s *Server) GetSignaturesForAddress(ctx context.Context, req *GetSignaturesForAddressRequest) (*GetSignaturesForAddressResponse, error) {
	opts := storage.SignatureOpts{
		Limit:      req.Limit,
		Before:     nil,
		Until:      nil,
		Commitment: storage.Commitment(req.Commitment),
	}
	if req.Before != "" {
		opts.Before = &req.Before
	}
	if req.Until != "" {
		opts.Until = &req.Until
	}
	sigs, err := s.root.GetSignaturesForAddress(ctx, req.Address, opts)
	if err != nil {
		return nil, err
	}
	out := make([]*SigInfo, len(sigs))
	for i, si := range sigs {
		out[i] = &SigInfo{
			Signature: si.Signature,
			Slot:      si.Slot,
			Err:       strPtrStr(si.Err),
			Memo:      strPtrStr(si.Memo),
			BlockTime: si.BlockTime.Unix(),
		}
	}
	return &GetSignaturesForAddressResponse{Signatures: out}, nil
}

func (s *Server) GetBlocksWithLimit(ctx context.Context, req *GetBlocksWithLimitRequest) (*GetBlocksWithLimitResponse, error) {
	slots, err := s.root.GetBlocksWithLimit(ctx, req.StartSlot, req.Limit, storage.Commitment(req.Commitment))
	if err != nil {
		return nil, err
	}
	return &GetBlocksWithLimitResponse{Slots: slots}, nil
}

func (s *Server) GetBlockTime(ctx context.Context, req *GetBlockTimeRequest) (*GetBlockTimeResponse, error) {
	t, err := s.root.GetBlockTime(ctx, req.Slot)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return &GetBlockTimeResponse{}, nil
	}
	return &GetBlockTimeResponse{BlockTime: t.Unix()}, nil
}

func strPtrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}