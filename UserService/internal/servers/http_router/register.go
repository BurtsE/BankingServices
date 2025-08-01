package http_router

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

func (r *Router) registerUserHandler(w http.ResponseWriter, req *http.Request) {

	propagator := otel.GetTextMapPropagator()
	reqCtx := propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

	ctx, span := r.tracer.Start(reqCtx, "register handler")
	defer span.End()

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

	user, err := r.service.Register(ctx, reqBody.Email, reqBody.Username, reqBody.Password, reqBody.FullName)
	if err != nil {
		r.logger.WithError(err).Error("failed to register user")
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
