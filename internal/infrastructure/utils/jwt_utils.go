// file: internal/infrastructure/utils/jwt_utils.go
package utils

import (
	"errors"
	"net/http"
	"strings"
)

// ExtractBearerToken extracts a Bearer token from the Authorization header.
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}