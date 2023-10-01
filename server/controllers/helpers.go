package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/session"
	"log"
	"net/http"
)

func ReturnMessageJSON(w http.ResponseWriter, message string, httpCode int, status string) {
	w.Header().Set("Content-Type", "application/json")
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

func CheckIfUserLoggedin(r *http.Request) bool {
	sessionData, err := session.SessionStorage.GetSession(r)
	if err != nil {
		log.Fatal(err)
	}
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
	SessionData, err := session.SessionStorage.GetSession(r)
	if err != nil {
		log.Fatal(err)
	}
	userID := SessionData.UserId
	username, _ := database.GetUsername(userID)
	return username
}

func getUserId(r *http.Request) int {
	SessionData, err := session.SessionStorage.GetSession(r)
	if err != nil {
		log.Fatal(err)
	}
	userID := SessionData.UserId
	return userID
}

func RateLimited(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ReturnMessageJSON(w, "Rate OK", http.StatusOK, "success")
}
