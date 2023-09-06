package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/session"
	"net/http"
)

func returnMessageJSON(w http.ResponseWriter, message string, httpCode int, status string) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(database.Response{
		Status:  status,
		Message: message,
	})
}

func checkForSessionToken(r *http.Request) (*http.Cookie, bool) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil, false
	}

	return sessionToken, true
}

func checkIfUserLoggedin(sessionToken *http.Cookie) bool {
	sessionData := session.SessionStorage.GetSession(sessionToken.Value)
	return sessionData.UserId != 0
}

func containsStr(arr []string, str rune) bool {
	for _, word := range arr {
		for _, letter := range word {
			if letter == str {
				return false
			}
		}
	}
	return true
}
