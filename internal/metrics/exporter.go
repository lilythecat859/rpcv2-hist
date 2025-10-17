package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func ServePrometheus(addr string, log *zap.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Info("starting prometheus metrics", zap.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("metrics server failed", zap.Error(err))
	}
}