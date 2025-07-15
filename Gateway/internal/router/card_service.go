package router

import (
	"net/http"
	"time"
)

const CARD_LABEL = "card"

func (r *Router) CardServiceHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		r.metrics.Duration.WithLabelValues(CARD_LABEL).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(CARD_LABEL).Inc()

	_, err := r.getIDFromJWTHeader(req)
	if err != nil {
		r.logger.Debugf("Error getting ID from JWT: %v", err)
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	r.proxy.ServeHTTP(w, req)
}
