package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type transferRequest struct {
	UUID          string `json:"uuid"`
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
}

func (r *Router) transferHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody transferRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if reqBody.FromAccountID == reqBody.ToAccountID {
		http.Error(w, "Invalid transfer parameters", http.StatusBadRequest)
		return
	}

	account, err := r.service.GetAccountByID(context.Background(), reqBody.FromAccountID)
	if err != nil || account.UserID != reqBody.UUID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "Could not get account", http.StatusBadRequest)
		return
	}
	if err = r.service.Transfer(req.Context(), reqBody.FromAccountID, reqBody.ToAccountID, reqBody.Amount); err != nil {
		r.logger.WithError(err).Error("Transfer failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
