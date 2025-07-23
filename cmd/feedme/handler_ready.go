package main

import (
	"net/http"

	"github.com/mahmoudk1000/feedme/pkg/utils"
)

func handlerReady(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}
