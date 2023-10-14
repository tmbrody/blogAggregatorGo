package main

import (
	"net/http"
	"strconv"

	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func (apiCfg *apiConfig) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, _ := ctx.Value(dbContextKey).(*database.Queries)

	limitString := r.URL.Query().Get("limit")

	if limitString == "" {
		limitString = "50"
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid limit")
		return
	}

	args := database.GetPostsByUserParams{
		UserID: apiCfg.User.ID,
		Limit:  int32(limit),
	}

	post, err := db.GetPostsByUser(ctx, args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get posts for the user")
		return
	}

	respondWithJSON(w, http.StatusOK, post)
}
