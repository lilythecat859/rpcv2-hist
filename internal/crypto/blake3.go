//go:build !cgo
// +build !cgo

package crypto

import (
	"encoding/binary"

	"github.com/zeebo/blake3"
)

// Hash64 returns the first 8 bytes of Blake3 hash as uint64 for fast deterministic sharding.
func Hash64(b []byte) uint64 {
	h := blake3.Sum256(b)
	return binary.LittleEndian.Uint64(h[:8])
}