package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type depositRequest struct {
	UUID      string `json:"uuid"`
	AccountID int64  `json:"account_id"`
	Amount    string `json:"amount"`
}

func (r *Router) depositHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody depositRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	account, err := r.service.GetAccountByID(context.Background(), reqBody.AccountID)
	if err != nil || account.UserID != reqBody.UUID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "could not get account", http.StatusInternalServerError)
		return
	}

	if err = r.service.Deposit(req.Context(), account.ID, reqBody.Amount); err != nil {
		r.logger.WithError(err).Error("failed to deposit")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
