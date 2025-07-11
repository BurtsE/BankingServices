package router

import (
	"errors"
	"net/http"
	"strings"
)

func extractJWTFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenStr, nil
}
