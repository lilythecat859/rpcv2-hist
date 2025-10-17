package model

import (
	"encoding/json"
	"time"
)

// Block represents a Solana slot/block.
type Block struct {
	Slot          uint64         `json:"slot" ch:"slot"`
	Blockhash     string         `json:"blockhash" ch:"blockhash"`
	ParentSlot    uint64         `json:"parentSlot" ch:"parent_slot"`
	BlockTime     int64          `json:"blockTime" ch:"block_time"` // unix seconds
	Height        uint64         `json:"blockHeight" ch:"height"`
	TxSigs        []string       `json:"transactions,omitempty" ch:"-"`
	Txs           []Transaction  `json:"-" ch:"-"`
	Raw           json.RawMessage `json:"-" ch:"raw"` // compressed raw block
}

// Transaction represents a solana transaction.
type Transaction struct {
	Signature       string          `json:"signature" ch:"signature"`
	Slot            uint64          `json:"slot" ch:"slot"`
	Index           uint64          `json:"-" ch:"tx_idx"` // position inside block
	BlockTime       int64           `json:"blockTime" ch:"block_time"`
	Signer          string          `json:"signer" ch:"signer"` // first signer
	Fee             uint64          `json:"fee" ch:"fee"`
	ComputeUnits    uint64          `json:"computeConsumed" ch:"compute_units"`
	Err             *string         `json:"err,omitempty" ch:"err"`
	Raw             json.RawMessage `json:"-" ch:"raw"` // compressed raw tx
}

// SignatureInfo is a lightweight row returned by getSignaturesForAddress.
type SignatureInfo struct {
	Signature string    `json:"signature"`
	Slot      uint64    `json:"slot"`
	Err       *string   `json:"err,omitempty"`
	Memo      *string   `json:"memo,omitempty"`
	BlockTime time.Time `json:"blockTime"`
}