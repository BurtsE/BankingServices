package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type withdrawRequest struct {
	UUID      string `json:"uuid"`
	AccountID int64  `json:"account_id"`
	Amount    string `json:"amount"`
}

func (r *Router) withdrawHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody withdrawRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("withdraw: json decode error: %s", err)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UUID == "" || reqBody.Amount == "" || reqBody.AccountID == 0 {
		r.logger.Debugf("withdraw: missing parameters: %v", reqBody)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	account, err := r.service.GetAccountByID(context.Background(), reqBody.AccountID)
	if err != nil || account.UserID != reqBody.UUID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "could not get account", http.StatusInternalServerError)
		return
	}

	if err = r.service.Withdraw(req.Context(), account.ID, reqBody.Amount); err != nil {
		r.logger.WithError(err).Error("failed to withdraw")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
