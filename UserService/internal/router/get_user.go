package router

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (r *Router) getUserByIDHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	user, err := r.userService.GetByID(ctx, userID)
	if err != nil {
		r.logger.WithError(err).Warn("user not found")
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
