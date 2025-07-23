package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
	"github.com/mahmoudk1000/feedme/internal/models"
	"github.com/mahmoudk1000/feedme/pkg/utils"
)

func (cfg *apiConfig) userCreateHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := data{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func (cfg *apiConfig) userGetHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	utils.RespondWithJson(w, http.StatusOK, models.DatabaseUserToUser(user))
}
