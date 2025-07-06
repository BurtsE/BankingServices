package router

import (
	"encoding/json"
	"net/http"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

func (r *Router) registerUserHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody registerRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	user, err := r.userService.Register(ctx, reqBody.Email, reqBody.Username, reqBody.Password, reqBody.FullName)
	if err != nil {
		r.logger.WithError(err).Error("failed to register user")
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
