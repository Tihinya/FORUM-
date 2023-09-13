package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/session"
	"net/http"
)

func ReturnMessageJSON(w http.ResponseWriter, message string, httpCode int, status string) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(database.Response{
		Status:  status,
		Message: message,
	})
}

func CheckForSessionToken(r *http.Request) (*http.Cookie, bool) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil, false
	}

	return sessionToken, true
}

func CheckIfUserLoggedin(sessionToken *http.Cookie) bool {
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

func getUsername(r *http.Request) string {
	sessionToken, _ := CheckForSessionToken(r)
	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, _ := database.GetUsername(userID)
	return username
}

func getUserId(r *http.Request) int {
	sessionToken, _ := CheckForSessionToken(r)
	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	return userID
}
