//go:build !cgo
// +build !cgo

package compress

import (
	"bytes"
	"fmt"

	"github.com/pierrec/lz4/v4"
)

// Compress returns lz4 compressed bytes.
func Compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := lz4.NewWriter(&buf)
	if _, err := w.Write(src); err != nil {
		return nil, fmt.Errorf("lz4 write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("lz4 close: %w", err)
	}
	return buf.Bytes(), nil
}

// Decompress lz4 bytes into dst (resize if needed).
func Decompress(src []byte) ([]byte, error) {
	r := lz4.NewReader(bytes.NewReader(src))
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("lz4 read: %w", err)
	}
	return buf.Bytes(), nil
}