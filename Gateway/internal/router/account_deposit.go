package router

import (
	"bytes"
	"encoding/json"
	"gateway/internal/domain"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const DepositAccountLabel = "User registry"

type depositAccountRequest struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
	Amount    string `json:"amount"`
}

func (r *Router) AccountDepositHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		r.metrics.Duration.WithLabelValues(CreateAccountLabel).Observe(time.Since(start).Seconds())
	}()
	r.metrics.Requests.WithLabelValues(CreateAccountLabel).Inc()

	var reqBody depositAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("cannot decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.UserID == "" || reqBody.AccountID == "" || reqBody.Amount == "" {
		r.logger.Debugf("insufficient request body: %v", reqBody)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	depositInfo := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(depositInfo).Encode(reqBody)

	event, err := domain.NewEventRequest(uuid.New(), time.Now(), depositInfo.Bytes(), "deposit_account")
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
