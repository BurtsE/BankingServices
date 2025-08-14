package router

import (
	"bytes"
	"encoding/json"
	"gateway/internal/domain"
	"net/http"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

func (r *Router) RegisterUserHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody registerRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("cannot decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.Email == "" || reqBody.Username == "" || reqBody.Password == "" || reqBody.FullName == "" {
		r.logger.Debugf("insufficient request body: %v", reqBody)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userInfo := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(userInfo).Encode(reqBody)

	event, err := domain.NewEventRequest(userInfo.Bytes(), "create_user")
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	r.requestService.AddEvent(req.Context(), event)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponceBody{EventId: event.ID().String()})
}
