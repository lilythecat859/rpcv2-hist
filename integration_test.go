//go:build integration
// +build integration

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/faithful-rpc/rpcv2-hist/internal/api/jsonrpc"
	"github.com/faithful-rpc/rpcv2-hist/internal/fractal"
	"github.com/faithful-rpc/rpcv2-hist/internal/model"
	"github.com/faithful-rpc/rpcv2-hist/internal/storage/clickhouse"
	"github.com/faithful-rpc/rpcv2-hist/internal/testutil"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	conn := testutil.ClickHouseTestContainer(t)
	db := &clickhouse.DB{Conn: conn}
	root := fractal.NewRoot(db, nil)
	srv := jsonrpc.NewServer(root, nil)

	block := &model.Block{
		Slot:      42,
		Blockhash: "fake",
		Raw:       json.RawMessage(`{"foo":"bar"}`),
	}
	require.NoError(t, conn.Exec(ctx, `INSERT INTO blocks (slot,blockhash,parent_slot,block_time,height,commitment,raw) VALUES (42,'fake',0,0,0,'finalized','')`))

	req := jsonrpc.request{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "getBlock",
		Params:  json.RawMessage(`[42]`),
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest("POST", "/", mustMarshal(req)))
	require.Equal(t, http.StatusOK, rec.Code)
	var resp jsonrpc.response
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Nil(t, resp.Error)
	require.NotNil(t, resp.Result)
}

func mustMarshal(v interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		panic(err)
	}
	return &buf
}