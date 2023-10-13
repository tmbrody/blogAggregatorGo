package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func (apiCfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	var params struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	feeds, err := db.GetAllFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get all feeds")
		return
	}

	feedID := uuid.UUID{}
	for _, feed := range feeds {
		if params.FeedID == feed.ID.String() {
			feedID = feed.ID
			break
		}
	}

	feedFollowID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate feed_follow ID")
		return
	}

	args := database.CreateFeedFollowParams{
		ID:        feedFollowID,
		FeedID:    feedID,
		UserID:    apiCfg.User.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	feedFollow, err := db.CreateFeedFollow(ctx, args)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to create new feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollow)
}

func (apiCfg *apiConfig) getUserFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	feedFollows, err := db.GetAllFeedFollows(ctx, apiCfg.User.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get all feed follows for user")
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (apiCfg *apiConfig) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	feedFollowString := chi.URLParam(r, "feedFollowID")

	feedFollows, err := db.GetAllFeedFollows(ctx, apiCfg.User.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to get all feed follows for user")
		return
	}

	feedFollowID := uuid.UUID{}
	for _, feedFollow := range feedFollows {
		if feedFollow.ID.String() == feedFollowString {
			feedFollowID = feedFollow.ID
			break
		}
	}

	args := database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: apiCfg.User.ID,
	}

	_, err = db.DeleteFeedFollow(ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, "Feed follow deleted")
}
