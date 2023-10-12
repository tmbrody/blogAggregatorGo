package main

import "net/http"

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && (authHeader[:7] == "Bearer " || authHeader[:7] == "ApiKey ") {
		return authHeader[7:]
	}
	return ""
}
