package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type showCardsRequest struct {
	AccountID string `json:"account_id"`
}

func (r *Router) showCardsHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody showCardsRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("show cards: json decode error: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if reqBody.AccountID == "" {
		r.logger.Debugf("show cards: not enough parameters: %v", reqBody)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), DefaultTimeout)
	defer cancel()

	isActive, err := r.banking.AccountIsActive(ctx, reqBody.AccountID)
	if err != nil {
		r.logger.WithError(err).Error("could not validate account")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isActive {
		r.logger.WithError(err).Error("account deactivated")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := r.service.GetCardsByAccount(ctx, reqBody.AccountID)
	if err != nil {
		r.logger.WithError(err).Error("could not get cards")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
