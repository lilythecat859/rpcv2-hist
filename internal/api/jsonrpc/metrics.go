package jsonrpc

import (
	"github.com/faithful-rpc/rpcv2-hist/internal/telemetry"
)

var metrics = telemetry.NewMetrics()

func record(method string, status string, dur float64) {
	metrics.RecordRequest(method, status, dur)
}