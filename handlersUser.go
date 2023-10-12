package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	var params struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate user ID")
		return
	}

	args := database.CreateUserParams{
		ID:        userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	}

	user, err := db.CreateUser(ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create new user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := extractTokenFromHeader(r)
	if apiKey == "" {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid or missing")
		return
	}

	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	user, err := db.GetUserByApiKey(ctx, apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to find user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
