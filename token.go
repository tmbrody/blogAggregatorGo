package main

import (
	"net/http"
	"strings"
)

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && strings.Contains(authHeader, " ") {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}
