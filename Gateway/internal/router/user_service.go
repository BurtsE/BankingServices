package router

import "net/http"

func (r *Router) UserServiceHandler(w http.ResponseWriter, req *http.Request) {
	r.metrics.Requests.WithLabelValues("user").Inc()

	req.Header.Set("X-Forwarded-For", req.RemoteAddr)
	r.proxy.ServeHTTP(w, req)
}
