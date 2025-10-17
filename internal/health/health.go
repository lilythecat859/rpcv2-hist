package health

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	store  Pingable
	log    *zap.Logger
	tracer trace.Tracer
}

type Pingable interface {
	Ping(context.Context) error
}

func NewHandler(store Pingable, log *zap.Logger) http.Handler {
	return &Handler{
		store:  store,
		log:    log,
		tracer: otel.Tracer("health"),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "health")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	status := http.StatusOK
	if err := h.store.Ping(ctx); err != nil {
		h.log.Warn("health ping fail", zap.Error(err))
		status = http.StatusServiceUnavailable
	}
	w.WriteHeader(status)
}