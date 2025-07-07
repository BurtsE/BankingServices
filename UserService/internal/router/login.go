package router

import (
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *Router) loginHandler(w http.ResponseWriter, req *http.Request) {
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

	ctx := req.Context()
	jwtToken, expiresAt, err := r.service.Authenticate(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		r.logger.WithError(err).Warn("failed to authenticate user")
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      jwtToken,
		"expires_at": expiresAt,
	})
}
