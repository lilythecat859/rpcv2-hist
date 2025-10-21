package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"go.uber.org/zap"

	"github.com/lilythecat859/rpcv2-hist/internal/model"
	"github.com/lilythecat859/rpcv2-hist/internal/storage"
)

type Ingester struct {
	store  storage.HistoricalStore
	log    *zap.Logger
	tick   time.Duration
	queue  chan *batch
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

type batch struct {
	blocks []model.Block
	txs    []model.Transaction
	sigs   []sigRow
}

type sigRow struct {
	address   string
	signature string
	slot      uint64
}

type Option func(*Ingester)

func WithLogger(l *zap.Logger) Option {
	return func(i *Ingester) { i.log = l }
}

func New(store storage.HistoricalStore, opts ...Option) (*Ingester, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ing := &Ingester{
		store:  store,
		log:    zap.NewNop(),
		tick:   400 * time.Millisecond,
		queue:  make(chan *batch, 1024),
		ctx:    ctx,
		cancel: cancel,
	}
	for _, o := range opts {
		o(ing)
	}
	return ing, nil
}

func (i *Ingester) Run(ctx context.Context) error {
	i.wg.Add(1)
	go i.loop()
	<-ctx.Done()
	i.cancel()
	i.wg.Wait()
	return nil
}

func (i *Ingester) loop() {
	defer i.wg.Done()
	ticker := time.NewTicker(i.tick)
	defer ticker.Stop()
	for {
		select {
		case <-i.ctx.Done():
			return
		case b := <-i.queue:
			if err := i.flush(b); err != nil {
				i.log.Error("flush batch", zap.Error(err))
			}
		case <-ticker.C:
		}
	}
}

func (i *Ingester) EnqueueBlock(block *model.Block) {
	i.queue <- &batch{blocks: []model.Block{*block}}
}

func (i *Ingester) flush(b *batch) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return i.storeBulk(ctx, b)
}

func (i *Ingester) storeBulk(ctx context.Context, b *batch) error {
	pool := memory.NewGoAllocator()
	defer pool.AssertSize(0)

	const maxRows = 1 << 20
	bldBlocks := array.NewStructBuilder(pool, arrow.StructOf(
		arrow.Field{Name: "slot", Type: arrow.PrimitiveTypes.Uint64},
		arrow.Field{Name: "blockhash", Type: arrow.BinaryTypes.String},
		arrow.Field{Name: "parent_slot", Type: arrow.PrimitiveTypes.Uint64},
		arrow.Field{Name: "block_time", Type: arrow.PrimitiveTypes.Int64},
		arrow.Field{Name: "height", Type: arrow.PrimitiveTypes.Uint64},
		arrow.Field{Name: "raw", Type: arrow.BinaryTypes.Binary},
	))
	defer bldBlocks.Release()

	for _, blk := range b.blocks {
		bldBlocks.Append(true)
		bldBlocks.FieldBuilder(0).(*array.Uint64Builder).Append(blk.Slot)
		bldBlocks.FieldBuilder(1).(*array.StringBuilder).Append(blk.Blockhash)
		bldBlocks.FieldBuilder(2).(*array.Uint64Builder).Append(blk.ParentSlot)
		bldBlocks.FieldBuilder(3).(*array.Int64Builder).Append(blk.BlockTime)
		bldBlocks.FieldBuilder(4).(*array.Uint64Builder).Append(blk.Height)
		bldBlocks.FieldBuilder(5).(*array.BinaryBuilder).Append(blk.Raw)
	}
	return nil
}