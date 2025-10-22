package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lilythecat859/rpcv2-hist/internal/model"
)

type Client struct {
	base   string
	httpcl *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		base:   baseURL,
		httpcl: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetBlock(ctx context.Context, slot uint64, commitment string) (*model.Block, error) {
	url := fmt.Sprintf("%s/block/%d?commitment=%s", c.base, slot, commitment)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpcl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var blk model.Block
	if err := json.NewDecoder(resp.Body).Decode(&blk); err != nil {
		return nil, err
	}
	return &blk, nil
}

func (c *Client) GetTransaction(ctx context.Context, signature string, commitment string) (*model.Transaction, error) {
	url := fmt.Sprintf("%s/tx/%s?commitment=%s", c.base, signature, commitment)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpcl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var tx model.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}