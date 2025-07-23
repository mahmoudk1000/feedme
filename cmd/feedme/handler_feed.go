package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
	"github.com/mahmoudk1000/feedme/internal/models"
	"github.com/mahmoudk1000/feedme/pkg/utils"
)

func (cfg *apiConfig) createFeedHandler(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type data struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := data{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create feed")
		return
	}

	feedFollow, err := cfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		log.Println("Error following feed:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, struct{
		Feed database.Feed `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
			Feed: feed,
			FeedFollow: feedFollow,
		})
}

func (cfg *apiConfig) getAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Printf("error: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if len(feeds) < 1 {
		utils.RespondWithError(w, http.StatusNotFound, "feeds are empty")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, models.DatabaseFeedsToFeeds(feeds))
}
