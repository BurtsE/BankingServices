package router

import (
	"encoding/json"
	"net/http"
)

type createAccountRequest struct {
	UserID         string `json:"user_id"`
	AccountType    string `json:"account_type"`
	AccountSubType string `json:"account_subtype"`
	Currency       string `json:"currency"`
}

func (r *Router) createAccountHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody createAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("createAccount: json decode error: %s", err)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UserID == "" || reqBody.Currency == "" ||
		reqBody.AccountType == "" || reqBody.AccountSubType == "" {
		r.logger.Debugf("createAccount: missing parameters: %v", reqBody)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	account, err := r.service.CreateAccount(req.Context(),
		reqBody.UserID,
		reqBody.Currency,
		reqBody.AccountType,
		reqBody.AccountSubType,
	)
	if err != nil {
		r.logger.WithError(err).Error("account creation fail")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}
