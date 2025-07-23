package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: apikey <api_key>
func GetApiKey(h http.Header) (string, error) {
	authHeader := h.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no apikey was provided")
	}

	apiKey := strings.Split(authHeader, " ")
	if len(apiKey) < 2 {
		return "", errors.New("malformed authorization header")
	}

	if apiKey[0] != "apikey" {
		return "", errors.New("malformed authorization header type")
	}

	return apiKey[1], nil
}
