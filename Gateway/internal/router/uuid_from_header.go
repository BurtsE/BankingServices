package router

import (
	"errors"
	"net/http"
	"strings"
)

func (r *Router) getIDFromJWTHeader(req *http.Request) (string, error) {
	jwt, err := extractJWTFromHeader(req)
	if err != nil {
		r.logger.Debugf("Error extracting JWT from header: %v", err)
		return "", err
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
		return "", err
	}

	if !ok {
		r.logger.Debugf("No token in cache: %v", jwt)
		uuid, err = r.userService.Validate(jwt)
		if err != nil {
			r.logger.Errorf("Error validating token: %v", err)
			return "", err
		}

		r.logger.Debugf("caching token: %v", jwt)
		err = r.tokenCache.Save(req.Context(), jwt, uuid, cachingDuration)
		if err != nil {
			r.logger.Errorf("Error saving token: %v", err)
		}
	}

	return uuid, nil
}

func extractJWTFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenStr, nil
}
