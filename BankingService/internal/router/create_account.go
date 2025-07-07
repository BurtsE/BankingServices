package router

import (
	"encoding/json"
	"net/http"
)

type createAccountRequest struct {
	UUID     string `json:"uuid"`
	Currency string `json:"currency"`
}

func (r *Router) createAccountHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody createAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("createAccount: json decode error: %s", err)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.UUID == "" || reqBody.Currency == "" {
		r.logger.Debugf("createAccount: missing parameters: %v", reqBody)
		http.Error(w, errInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	switch reqBody.Currency {
	case "usd":
	case "rub":
	default:
		http.Error(w, "Invalid currency", http.StatusBadRequest)
		return
	}

	account, err := r.service.CreateAccount(req.Context(), reqBody.UUID, reqBody.Currency)
	if err != nil {
		r.logger.WithError(err).Error("account creation fail")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}
