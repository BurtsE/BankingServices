package router

import (
	"net/http"
)

func (r *Router) BankingServiceHandler(w http.ResponseWriter, req *http.Request) {
	r.metrics.Requests.WithLabelValues("banking").Inc()

	jwt, err := extractJWTFromHeader(req)
	if err != nil {
		r.logger.Debugf("Error extracting JWT from header: %v", err)
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var (
		uuid string
		ok   bool
	)

	if uuid, ok, err = r.tokenCache.Get(req.Context(), jwt); err != nil {
		r.logger.Errorf("Error getting token from cache: %v", err)
	}
	if err != nil {
		r.logger.Debugf("Error getting token from cache: %v", err)
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	if !ok {
		r.logger.Debugf("No token in cache: %v", jwt)
		uuid, err = r.userService.Validate(jwt)
		if err != nil {
			r.logger.Errorf("Error validating token: %v", err)
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		r.logger.Debugf("caching token: %v", jwt)
		err = r.tokenCache.Save(req.Context(), jwt, uuid, caching_duration)
		if err != nil {
			r.logger.Debugf("Error saving token: %v", err)
		}
	}

	err = insertUserIDToRequestBody(req, uuid)
	if err != nil {
		r.logger.Debugf("Error saving token: %v", err)
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusInternalServerError)
		return
	}

	r.proxy.ServeHTTP(w, req)
}
