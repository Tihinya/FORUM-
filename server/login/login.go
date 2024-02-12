package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/config"
	"forum/database"
	"forum/security"
	"forum/validation"
	"io"
	"log"
	"net/http"
)

type githubUser struct {
	Login string `json:"login"`
	Email string `json:"email"`
}
type userData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetGithubAccessToken(code string) (string, error) {

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     config.Config.GitHubClientId,
		"client_secret": config.Config.GitHubClientSecret,
		"code":          code,
		"scope":         "user:email",
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return "", reqerr
	}

	// Set return type JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken, nil
}

func GetGithubData(accessToken string) (githubUser, error) {

	// Parsing to workable struct
	var user githubUser

	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)

	if reqerr != nil {
		return githubUser{}, reqerr
	}

	// Set the Authorization header before sending the request
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)
	req.Header.Set("Accept", "application/vnd.github+json")

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return githubUser{}, resperr
	}

	// Read the response as a byte slice
	respbody, _ := io.ReadAll(resp.Body)

	err := json.Unmarshal(respbody, &user)
	if err != nil {
		return githubUser{}, err
	}

	// Convert byte slice to string and return
	return user, nil
}
func AddUser(username string, email string, password string, roleId int, age, gender string) (int, error) {
	user := database.UserInfo{
		ProfilePicture: config.Config.ProfilePicture,
		Username:       username,
		Email:          email,
		Password:       password,
		RoleID:         roleId,
		Age:            age,
		Gender:         gender,
	}

	// Encrypt the password
	hashedPassword, err := security.PasswordEncrypting(password)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)

	id, err := database.CreateUser(user)
	if err != nil {
		return 0, err
	}
	user.ID = id
	return id, nil
}
func GetUserData(token string) userData {
	var result userData
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		log.Println(err)
		return result
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return result
	}
	var userInfo userData
	if err := json.Unmarshal(body, &userInfo); err != nil {
		log.Println(err)
		return result
	}
	result = userInfo
	return result
}

func CreateAdminUser() {
	// Check if admin user already exists
	roleId, err := database.GetRoleId("admin")
	if err != nil {
		log.Println(err)
	}

	id, err := validation.GetUserID(database.DB, "", "admin")
	if err != nil {
		log.Println(err)
	}
	if id == 0 {
		// User does not exist, create a new user
		_, err = AddUser("admin", "admin@example.com", "admin", roleId, "17", "male")
		if err != nil {
			log.Println(err)
		}
	}
}
