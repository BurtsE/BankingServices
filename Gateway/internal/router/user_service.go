package router

import (
	"net/http"
	"time"
)

const USER_LABEL = "user"

func (r *Router) UserServiceHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		r.metrics.Duration.WithLabelValues(USER_LABEL).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(USER_LABEL).Inc()

	req.Header.Set("X-Forwarded-For", req.RemoteAddr)
	r.proxy.ServeHTTP(w, req)
}
