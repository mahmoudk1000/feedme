package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
	"github.com/mahmoudk1000/feedme/internal/models"
	"github.com/mahmoudk1000/feedme/pkg/utils"
)

func (cfg *apiConfig) followFeedHanlder(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type data struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := data{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	isAlreadyFollowing, err := cfg.DB.CheckFeedFollow(r.Context(), database.CheckFeedFollowParams{
		FeedID: params.FeedID,
		UserID: user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check feed follow")
		return
	}
	if isAlreadyFollowing {
		utils.RespondWithError(w, http.StatusFound, "feed already followed")
		return
	}

	feedFollow, err := cfg.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
	})
	if err != nil {
		log.Println("Error following feed:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) getUserFeedFollowes(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	userFeeds, err := cfg.DB.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		log.Println("Error getting user feeds:", err)
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("error getting %v 's feeds", user.Name),
		)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(userFeeds))
}

func (cfg *apiConfig) deleteFeedFollowHandler(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedIdStr := chi.URLParam(r, "feedFollowID")
	feedId, err := uuid.Parse(feedIdStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid feed ID")
		return
	}
	isAlreadyFollowing, err := cfg.DB.CheckFeedFollow(r.Context(), database.CheckFeedFollowParams{
		FeedID: feedId,
		UserID: user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check feed follow")
		return
	}
	if !isAlreadyFollowing {
		utils.RespondWithError(w, http.StatusNotFound, "feed follow not found")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete feed follow")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, struct{}{})
}
