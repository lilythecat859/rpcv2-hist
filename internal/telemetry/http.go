package telemetry

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func HTTPMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http",
		otelhttp.WithPublicEndpoint(),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}