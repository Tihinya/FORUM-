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
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var login database.UserInfo

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		returnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	login.Email = strings.TrimSpace(login.Email)
	login.Password = strings.TrimSpace(login.Password)

	userId, err := validation.GetUserIdFromEmail(database.DB, login.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userId == 0 {
		returnMessageJSON(w, "incorrect login or password", http.StatusBadRequest, "error")
		return
	}

	userPassword, err := database.GetUserPassword(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the entered password matches the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(userPassword.Password), []byte(login.Password))
	if err != nil {
		returnMessageJSON(w, "incorrect login or password", http.StatusBadRequest, "error")
		return
	}

	token := session.SessionStorage.CreateSession(userId)
	session.SessionStorage.SetCookie(token, w)

	json.NewEncoder(w).Encode(database.LoginResponse{
		Status:  "success",
		Message: "User logined successfully",
		ID:      userId,
	})
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

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var scope = config.Config.GoogleOAuth
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code",
		config.Config.GoogleID, config.Config.GoogleRedirectURI, scope)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found from callback request!", http.StatusInternalServerError)
		return
	}

	values := url.Values{
		"code":          {code},
		"client_id":     {config.Config.GoogleID},
		"client_secret": {config.Config.GoogleClientSecret},
		"redirect_uri":  {config.Config.GoogleRedirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm(config.Config.GoogleGetToken, values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type oauthToken struct {
		AccessToken string `json:"access_token"`
	}

	var authenticateToken oauthToken
	if err := json.Unmarshal(body, &authenticateToken); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	googleUser := login.GetUserData(authenticateToken.AccessToken)
	if googleUser.Name == "" {
		http.Error(w, "Failed to get user data from Google API", http.StatusInternalServerError)
		return
	}

	id, err := validation.GetUserIdFromEmail(database.DB, googleUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		// User does not exist, create a new user
		id, err = login.AddUser(googleUser.Name, googleUser.Email, security.CreateRandomPassword(10))
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

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		config.Config.GitHubClientId,
		config.Config.GitHubRedirectURI,
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
	id, err := validation.GetUserIdFromEmail(database.DB, githubData.Email)
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
