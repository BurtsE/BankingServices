package router

import (
	"context"
	"encoding/json"
	"net/http"
)

type blockCardRequest struct {
	AccountID string `json:"account_id"`
	Pan       string `json:"pan"`
}

func (r *Router) blockCardHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody blockCardRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Debugf("createCard: json decode error: %s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if reqBody.AccountID == "" || len(reqBody.Pan) != 16 {
		r.logger.Debugf("block card: wrong parameters: %v", reqBody)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), DefaultTimeout)
	defer cancel()

	blocked, err := r.service.BlockCard(ctx, reqBody.AccountID, reqBody.Pan)
	if err != nil {
		r.logger.WithError(err).Errorf("block card: ")
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}
	if !blocked {
		r.logger.Infof("could not block card: %s", reqBody.Pan)
		http.Error(w, "Operation failed", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
