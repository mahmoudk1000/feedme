package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("Error marshaling JSON:", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Server error:", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, errorResponse{
		Error: msg,
	})
}
