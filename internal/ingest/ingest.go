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

	"github.com/faithful-rpc/rpcv2-hist/internal/model"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage"
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