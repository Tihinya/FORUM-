package controllers

import (
	"encoding/json"
	"forum/database"
	"net/http"
)

func returnMessageJSON(w http.ResponseWriter, message string, httpCode int, status string) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(database.Response{
		Status:  status,
		Message: message,
	})
}
