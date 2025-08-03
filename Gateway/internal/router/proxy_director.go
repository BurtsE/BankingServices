package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (r *Router) director(req *http.Request) {
	route, _ := strings.CutPrefix(req.URL.Path, RouterPrefix)
	parts := strings.Split(route, "/")
	uri := r.mapping[parts[1]]

	if uri == nil {
		req.URL.Scheme = "invalid"
		req.URL.Host = "no-route"
	} else {
		req.URL.Scheme = "http"
		req.URL.Host = uri.Host
	}
	req.Header.Set("X-Forwarded-For", req.RemoteAddr)

	log.Printf("Forwarding request to %s | %s %s", req.URL.Host, req.Method, req.URL.Path)
}

func (r *Router) errorHandler(w http.ResponseWriter, req *http.Request, err error) {
	r.logger.Printf("Server error: %v", err)
	writeErrorResponse(w, err, http.StatusServiceUnavailable)
}

func (r *Router) modifyResponse(res *http.Response) error {
	res.Header.Set("gateway", "GoG")
	return nil
}

type errorResponse struct {
	Error string `json:"message"`
	Code  int    `json:"code"`
}

func writeErrorResponse(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{
		Error: err.Error(),
		Code:  code,
	})
}
