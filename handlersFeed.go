package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func (apiCfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	var params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	feedID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate feed ID")
		return
	}

	args := database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    apiCfg.User.ID,
	}

	feed, err := db.CreateFeed(ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create new feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, feed)
}

func (apiCfg *apiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	feeds, err := db.GetAllFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get all feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}
