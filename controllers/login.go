package controllers

import (
	"encoding/json"
	"fmt"
	"forum/config"
	"forum/database"
	"forum/login"
	"forum/security"
	"forum/session"
	"forum/validation"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// This line is assigning a fixed user ID for testing purposes.
	// You might want to modify this logic to authenticate the user properly.
	userID := 228 // test user ID
	login.AddLogin(w, userID)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	token, err := session.ValidateToken(r)
	if err != nil {
		// Handle the error appropriately (e.g., log or return an error response)
		return
	}
	s := session.SessionStorage.GetSession(token)
	s.RemoveSession()
	session.SessionStorage.DeleteCookie(w)
}

func GoogleLogin(w http.ResponseWriter, r *http.Request)    {}
func GoogleCallback(w http.ResponseWriter, r *http.Request) {}

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		config.Config.GitHubClientId,
		"https://localhost:8080/login/github/callback",
	)

	// Redirect the user to the GitHub page for authentication
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken, err := login.GetGithubAccessToken(code)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Bad request - problem with access tokens!",
		})
		return
	}

	githubData, err := login.GetGithubData(githubAccessToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Bad request - problem with user data!",
		})
		return
	}

	// Check if the user has made their GitHub email public
	if githubData.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Please make your GitHub email public to proceed with authentication",
		})
		return
	}

	// Check if the user with the GitHub email already exists in the database
	id, err := validation.GetUserID(database.DB, githubData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		// User does not exist, create a new user
		id, err = login.AddUser(githubData.Login, githubData.Email, security.CreateRandomPassword(10))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Set session
	token := session.SessionStorage.CreateSession(id)
	session.SessionStorage.SetCookie(token, w)

	redirectURL := "/login/github/redirect"
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func GithubCallbackRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are logged into the server")
}
