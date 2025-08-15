package router

import (
	"bytes"
	"encoding/json"
	"gateway/internal/domain"
	"net/http"
	"time"
)

const CreateAccountLabel = "User registry"

type createAccountRequest struct {
	UserID         string `json:"user_id"`
	AccountType    string `json:"account_type"`
	AccountSubType string `json:"account_subtype"`
	Currency       string `json:"currency"`
}

func (r *Router) AccountCreationHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		r.metrics.Duration.WithLabelValues(CreateAccountLabel).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(CreateAccountLabel).Inc()

	var reqBody createAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("cannot decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.UserID == "" || reqBody.AccountType == "" || reqBody.AccountSubType == "" || reqBody.Currency == "" {
		r.logger.Debugf("insufficient request body: %v", reqBody)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accountInfo := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(accountInfo).Encode(reqBody)

	event, err := domain.NewEventRequest(accountInfo.Bytes(), "create_account")
	if err != nil {
		http.Error(w, "Could not create account", http.StatusInternalServerError)
		return
	}

	err = r.requestService.AddEvent(req.Context(), event)
	if err != nil {
		r.logger.Errorf("could not process event: %v", err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponceBody{EventId: event.ID().String()})
}
