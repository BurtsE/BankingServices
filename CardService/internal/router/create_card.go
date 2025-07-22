package router

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type createCardRequest struct {
	AccountID      string `json:"account_id"`
	CardHolderName string `json:"card_holder_name"`
}

func (r *Router) createCardHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody createCardRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("createCard: json decode error: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if reqBody.AccountID == "" || reqBody.CardHolderName == "" {
		r.logger.Debugf("createCard: not enough parameters: %v", reqBody)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	isActive, err := r.banking.AccountIsActive(reqBody.AccountID)
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

	card, err := r.service.GenerateVirtualCard(ctx, reqBody.AccountID, reqBody.CardHolderName)
	if err != nil {
		r.logger.WithError(err).Error("card generation failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
