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
)

func Login(w http.ResponseWriter, r *http.Request) {
	userID := 1 // test user ID
	token := session.SessionStorage.CreateSession(userID)
	session.SessionStorage.SetCookie(token, w)
}

func LogOut(w http.ResponseWriter, r *http.Request) {

	s, err := session.SessionStorage.GetSession(r)
	if err != nil {
		log.Panicln(err)
	}
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
	roleId, err := database.GetRoleId("user")
	if err != nil {
		log.Println(err)
	}
	id, err := validation.GetUserID(database.DB, googleUser.Email, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		// User does not exist, create a new user
		id, err = login.AddUser(googleUser.Name, googleUser.Email, security.CreateRandomPassword(10), roleId)
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
	roleId, err := database.GetRoleId("user")
	if err != nil {
		log.Println(err)
	}

	// Check if the user with the GitHub email already exists in the database
	id, err := validation.GetUserID(database.DB, githubData.Email, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		// User does not exist, create a new user
		id, err = login.AddUser(githubData.Login, githubData.Email, security.CreateRandomPassword(10), roleId)
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
