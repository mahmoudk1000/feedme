package main

import (
	"fmt"
	"net/http"

	"github.com/mahmoudk1000/feedme/internal/auth"
	"github.com/mahmoudk1000/feedme/internal/database"
	"github.com/mahmoudk1000/feedme/pkg/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(h authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Sprintf("unauthorized: %v", err))
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "user not found")
			return
		}

		h(w, r, user)
	}
}
