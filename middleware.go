package main

import (
	"context"
	"net/http"

	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func withDB(next http.HandlerFunc, db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), dbContextKey, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (apiCfg *apiConfig) middlewareAuth(next http.HandlerFunc, db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := extractTokenFromHeader(r)
		if apiKey == "" {
			respondWithError(w, http.StatusUnauthorized, "API key is invalid or missing")
			return
		}

		ctx := context.WithValue(r.Context(), dbContextKey, db)

		user, err := db.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unable to find user")
			return
		}

		apiCfg.User = user

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
