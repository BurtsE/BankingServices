package router

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const UserLabel = "user"

func (r *Router) UserServiceHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	defer func() {
		r.metrics.Duration.WithLabelValues(UserLabel).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(UserLabel).Inc()

	req.Header.Set("X-Forwarded-For", req.RemoteAddr)

	// inject context values (trace id, etc.) in the headers
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(req.Context(), propagation.HeaderCarrier(req.Header))

	r.proxy.ServeHTTP(w, req)
}
