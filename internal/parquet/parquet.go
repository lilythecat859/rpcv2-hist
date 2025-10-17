package parquet

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/parquet"
	"github.com/apache/arrow/go/v15/parquet/compress"
	"github.com/apache/arrow/go/v15/parquet/file"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
	"go.uber.org/zap"

	"github.com/faithful-rpc/rpcv2-hist/internal/model"
)

const (
	// 128 MB row groups
	targetRowGroupSize = 128 * 1024 * 1024
)

type Writer struct {
	path   string
	logger *zap.Logger
}

type Option func(*Writer)

func WithLogger(l *zap.Logger) Option {
	return func(w *Writer) { w.logger = l }
}

// NewWriter returns a writer that flushes Arrow record batches to a parquet file.
func NewWriter(path string, opts ...Option) *Writer {
	w := &Writer{path: path, logger: zap.NewNop()}
	for _, o := range opts {
		o(w)
	}
	return w
}

// WriteBlocks appends blocks to a parquet file (creates if not exists).
func (w *Writer) WriteBlocks(blocks []model.Block) error {
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "slot", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
			{Name: "blockhash", Type: arrow.BinaryTypes.String, Nullable: false},
			{Name: "parent_slot", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
			{Name: "block_time", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
			{Name: "height", Type: arrow.PrimitiveTypes.Uint64, Nullable: false},
			{Name: "raw", Type: arrow.BinaryTypes.Binary, Nullable: false},
		}, nil,
	)

	bld := array.NewRecordBuilder(pool, schema)
	defer bld.Release()

	for _, blk := range blocks {
		bld.Field(0).(*array.Uint64Builder).Append(blk.Slot)
		bld.Field(1).(*array.StringBuilder).Append(blk.Blockhash)
		bld.Field(2).(*array.Uint64Builder).Append(blk.ParentSlot)
		bld.Field(3).(*array.Int64Builder).Append(blk.BlockTime)
		bld.Field(4).(*array.Uint64Builder).Append(blk.Height)
		bld.Field(5).(*array.BinaryBuilder).Append(blk.Raw)
	}
	rec := bld.NewRecord()
	defer rec.Release()

	// create parent dir
	if err := os.MkdirAll(filepath.Dir(w.path), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

// write parquet
	fw, err := file.CreateFileWriter(w.path,
		parquet.WithCompression(compress.Codecs.Lz4),
	)
	if err != nil {
		return fmt.Errorf("create writer: %w", err)
	}
	defer fw.Close()

	wr, err := pqarrow.NewFileWriter(schema, fw,
		pqarrow.WithRowGroupLength(targetRowGroupSize),
	)
	if err != nil {
		return fmt.Errorf("new arrow writer: %w", err)
	}
	defer wr.Close()

	if err := wr.Write(rec); err != nil {
		return fmt.Errorf("write record: %w", err)
	}
	return nil
}

// ReadBlocks reads a parquet file into memory.
func ReadBlocks(path string) ([]model.Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	rdr, err := file.OpenParquetFile(path, true)
	if err != nil {
		return nil, fmt.Errorf("open parquet: %w", err)
	}
	defer rdr.Close()

	ar, err := pqarrow.NewFileReader(rdr)
	if err != nil {
		return nil, fmt.Errorf("arrow reader: %w", err)
	}

	rec, err := ar.Read()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	defer rec.Release()

	var out []model.Block
	rows := int(rec.NumRows())
	for i := 0; i < rows; i++ {
		out = append(out, model.Block{
			Slot:       rec.Column(0).(*array.Uint64).Value(i),
			Blockhash:  rec.Column(1).(*array.String).Value(i),
			ParentSlot: rec.Column(2).(*array.Uint64).Value(i),
			BlockTime:  rec.Column(3).(*array.Int64).Value(i),
			Height:     rec.Column(4).(*array.Uint64).Value(i),
			Raw:        rec.Column(5).(*array.Binary).Value(i),
		})
	}
	return out, nil
}

// UploadParquet uploads a local parquet file to S3 (placeholder).
func UploadParquet(ctx context.Context, localPath, s3Key string) error {
	// TODO: integrate minio/aws-sdk-go-v2
	return nil
}