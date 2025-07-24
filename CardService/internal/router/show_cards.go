package router

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (r *Router) showCardsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	accountID := vars["accountID"]

	ctx, cancel := context.WithTimeout(req.Context(), DefaultTimeout)
	defer cancel()

	isActive, err := r.banking.AccountIsActive(ctx, accountID)
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

	cards, err := r.service.GetCardsByAccount(ctx, accountID)
	if err != nil {
		r.logger.WithError(err).Error("could not get cards")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
