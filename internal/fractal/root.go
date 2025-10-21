package fractal

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/lilythecat859/rpcv2-hist/internal/model"
	"github.com/lilythecat859/rpcv2-hist/internal/storage"
)

// Root is the top-level fractal node that fans out reads to shards.
type Root struct {
	store  storage.HistoricalStore
	shards []*shard
	log    *zap.Logger
	mu     sync.RWMutex
}

type shard struct {
	id    uint32
	store storage.HistoricalStore
}

func NewRoot(store storage.HistoricalStore, log *zap.Logger) *Root {
	r := &Root{
		store: store,
		log:   log,
	}
	// start with one shard; fractal split happens on load.
	r.shards = append(r.shards, &shard{id: 0, store: store})
	return r
}

func (r *Root) GetBlock(ctx context.Context, slot uint64, commitment storage.Commitment) (*model.Block, error) {
	sh := r.shardFor(slot)
	return sh.store.GetBlock(ctx, slot, commitment)
}

func (r *Root) GetBlocksWithLimit(ctx context.Context, start, limit uint64, commitment storage.Commitment) ([]uint64, error) {
	// naive: ask first shard; fractal depth can be added later.
	return r.shards[0].store.GetBlocksWithLimit(ctx, start, limit, commitment)
}

func (r *Root) GetBlockTime(ctx context.Context, slot uint64) (*time.Time, error) {
	sh := r.shardFor(slot)
	return sh.store.GetBlockTime(ctx, slot)
}

func (r *Root) GetTransaction(ctx context.Context, signature string, commitment storage.Commitment) (*model.Transaction, error) {
	slot := slotFromSignature(signature)
	sh := r.shardFor(slot)
	return sh.store.GetTransaction(ctx, signature, commitment)
}

func (r *Root) GetSignaturesForAddress(ctx context.Context, addr string, opts storage.SignatureOpts) ([]model.SignatureInfo, error) {
	// fan-out to all shards and merge; cache can be added later.
	var out []model.SignatureInfo
	for _, sh := range r.shards {
		part, err := sh.store.GetSignaturesForAddress(ctx, addr, opts)
		if err != nil {
			r.log.Warn("shard query failed", zap.Uint32("shard", sh.id), zap.Error(err))
			continue
		}
		out = append(out, part...)
	}
	return out, nil
}

func (r *Root) shardFor(slot uint64) *shard {
	// simple hash-split; can be replaced with fractal tree.
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.shards[slot%uint64(len(r.shards))]
}

func slotFromSignature(sig string) uint64 {
	// cheap deterministic hash of first 8 bytes of signature.
	var h uint64
	for i := 0; i < 8 && i < len(sig); i++ {
		h = h<<8 + uint64(sig[i])
	}
	return h
}