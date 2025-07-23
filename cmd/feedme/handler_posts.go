package main

import (
	"log"
	"net/http"

	"github.com/mahmoudk1000/feedme/internal/database"
	"github.com/mahmoudk1000/feedme/internal/models"
	"github.com/mahmoudk1000/feedme/pkg/utils"
)

func (cfg *apiConfig) getUserPostsHandler(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	posts, err := cfg.DB.GetPostsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("error getting posts for user %s: %v", user.ID, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get user posts")
		return
	}

	if len(posts) < 1 {
		utils.RespondWithError(w, http.StatusNotFound, "user has no posts")
	}

	utils.RespondWithJson(w, http.StatusOK, models.DatabasePostsToPosts(posts))
}
