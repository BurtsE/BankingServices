package router

import (
	"net/http"
	"time"
)

const BANKING_LABEL = "banking"

func (r *Router) BankingServiceHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		r.metrics.Duration.WithLabelValues(BANKING_LABEL).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(BANKING_LABEL).Inc()

	uuid, err := r.getIDFromJWTHeader(req)
	if err != nil {
		r.logger.Debugf("Error getting ID from JWT: %v", err)
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = insertUserIDToRequestBody(req, uuid)
	if err != nil {
		r.logger.Debugf("Error saving token: %v", err)
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusInternalServerError)
		return
	}

	r.proxy.ServeHTTP(w, req)
}
