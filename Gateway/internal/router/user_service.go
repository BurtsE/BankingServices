package router

import (
	"net/http"
	"time"
)

const UserLabel = "user"

func (r *Router) UserServiceHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		r.metrics.Duration.WithLabelValues(UserLabel).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(UserLabel).Inc()

	req.Header.Set("X-Forwarded-For", req.RemoteAddr)
	r.proxy.ServeHTTP(w, req)
}
