package router

import (
	"encoding/json"
	"net/http"
	"time"
)

const LoginUserLabel = "User login"

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *Router) UserLoginHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	defer func() {
		r.metrics.Duration.WithLabelValues(LoginUserLabel).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(LoginUserLabel).Inc()

	var reqBody loginRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("cannot decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		r.logger.Debugf("insufficient request body: %v", reqBody)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Send request to user service
}
