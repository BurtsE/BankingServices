package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type transferRequest struct {
	UserID        string `json:"user_id"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        string `json:"amount"`
}

func (r *Router) transferHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody transferRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("transfer: json decode error: %s", err)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserID == "" || reqBody.FromAccountID == "" ||
		reqBody.ToAccountID == "" || reqBody.Amount == "" {
		r.logger.Debugf("transfer: missing parameters: %v", reqBody)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.FromAccountID == reqBody.ToAccountID {
		http.Error(w, "Invalid transfer parameters", http.StatusBadRequest)
		return
	}

	account, err := r.service.GetAccountByID(context.Background(), reqBody.FromAccountID)
	if err != nil || account.UserID != reqBody.UserID {
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
