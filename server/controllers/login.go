package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/config"
	"forum/database"
	"forum/login"
	"forum/security"
	"forum/session"
	"forum/socket"
	"forum/validation"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var login database.UserInfo

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	login.Email = strings.TrimSpace(login.Email)
	login.Password = strings.TrimSpace(login.Password)

	userId, err := validation.GetUserID(database.DB, login.Email, "")
	if err != nil {
		log.Println("Failed getting user id:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userId == 0 {
		ReturnMessageJSON(w, "incorrect login or password", http.StatusBadRequest, "error")
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
		ReturnMessageJSON(w, "incorrect login or password", http.StatusBadRequest, "error")
		return
	}

	token := session.SessionStorage.CreateSession(userId)
	session.SessionStorage.SetCookie(token, w)

	ReturnMessageJSON(w, "User logined successfully", http.StatusOK, "success")
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	s, err := session.SessionStorage.GetSession(r)
	if err != nil {
		log.Panicln(err)
	}
	if s == nil {
		return
	}

	userId, err := session.GetUserId(r)
	if err != nil {
		ReturnMessageJSON(w, "Error fetching user id", http.StatusInternalServerError, "error")
		return
	}

	socket.Instance.RemoveClientWithId(userId)

	s.RemoveSession()
	session.SessionStorage.DeleteCookie(w)
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	scope := config.Config.GoogleOAuth
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code",
		config.Config.GoogleID, config.Config.GoogleRedirectURI, scope)
	json.NewEncoder(w).Encode(url)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		ReturnMessageJSON(w, "Code not found from callback request!", http.StatusInternalServerError, "error")
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
		ReturnMessageJSON(w, "Failed to get user data from Google API", http.StatusInternalServerError, "error")
		return
	}
	roleId, err := database.GetRoleId("user")
	if err != nil {
		log.Println(err)
	}
	id, err := validation.GetUserID(database.DB, googleUser.Email, "")
	if err != nil {
		ReturnMessageJSON(w, "User ID Problem", http.StatusInternalServerError, "error")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		// User does not exist, create a new user
		id, err = login.AddUser(googleUser.Name, googleUser.Email, security.CreateRandomPassword(10), roleId, "", "")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Set session
	token := session.SessionStorage.CreateSession(id)
	session.SessionStorage.SetCookie(token, w)

	http.Redirect(w, r, "https://localhost:3000/", http.StatusTemporaryRedirect)
}

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		config.Config.GitHubClientId,
		config.Config.GitHubRedirectURI,
	)

	// Redirect the user to the GitHub page for authentication
	json.NewEncoder(w).Encode(redirectURL)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken, err := login.GetGithubAccessToken(code)
	if err != nil {
		ReturnMessageJSON(w, "You are not authorized, problem with access token", http.StatusBadRequest, "error")
		return
	}

	githubData, err := login.GetGithubData(githubAccessToken)
	if err != nil {
		ReturnMessageJSON(w, "Problem with user data", http.StatusBadRequest, "error")
		return
	}

	// Check if the user has made their GitHub email public
	if githubData.Email == "" {
		ReturnMessageJSON(w, "Please make your GitHub email public to proceed with authenticationa", http.StatusBadRequest, "error")
		return
	}
	roleId, err := database.GetRoleId("user")
	if err != nil {
		log.Println(err)
	}

	// Check if the user with the GitHub email already exists in the database
	id, err := validation.GetUserID(database.DB, githubData.Email, "")
	if err != nil {
		ReturnMessageJSON(w, "GitHub email already exists", http.StatusInternalServerError, "error")
		return
	}

	if id == 0 {
		// User does not exist, create a new user
		id, err = login.AddUser(githubData.Login, githubData.Email, security.CreateRandomPassword(10), roleId, "", "")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Set session
	token := session.SessionStorage.CreateSession(id)
	session.SessionStorage.SetCookie(token, w)

	http.Redirect(w, r, "https://localhost:3000/", http.StatusTemporaryRedirect)
}

func CheckAuthorization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := getUsername(r)

	ReturnMessageJSON(w, "User is authenticated as "+username, http.StatusOK, "success")
}
