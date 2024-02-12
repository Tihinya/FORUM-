package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"

	"forum/database"
	"forum/session"
	"forum/validation"
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
		log.Println(err)
	}

	if sessionData == nil {
		return false
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
		log.Println(err)
	}
	userID := SessionData.UserId
	username, _ := database.GetUsername(userID)
	return username
}

func GetUserId(r *http.Request) int {
	SessionData, err := session.SessionStorage.GetSession(r)
	if err != nil {
		log.Println("Getting user ID failed: ", err)
	}

	userID := SessionData.UserId

	return userID
}

func RateLimited(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ReturnMessageJSON(w, "Rate OK", http.StatusOK, "success")
}

func randomUserData(userNumber int) database.UserInfo {
	genders := []string{"male", "female"}

	username := fmt.Sprintf("testUser%d", userNumber)
	email := fmt.Sprintf("testUser%d@example.com", userNumber)
	password := generateRandomString(25)
	passwordConfirmation := password
	gender := genders[rand.Intn(len(genders))]
	age := fmt.Sprint(rand.Intn(1000))

	return database.UserInfo{
		Username:             username,
		Email:                email,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
		Gender:               gender,
		Age:                  age,
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func compareTwoStructs[T any](obj1, obj2 T, ignore map[string]bool) error {
	value1 := reflect.ValueOf(obj1)
	value2 := reflect.ValueOf(obj2)

	if value1.Type().Name() != value2.Type().Name() {
		return fmt.Errorf("structs must have the same type")
	}

	for i := 0; i < value1.NumField(); i++ {
		name := value1.Type().Field(i).Name
		if _, exists := ignore[name]; exists {
			continue
		}

		fieldValue1 := value1.Field(i).Interface()
		fieldValue2 := value2.Field(i).Interface()

		if fieldValue1 != fieldValue2 {
			return fmt.Errorf("fields are not equal. %s: %v != %v", name, fieldValue1, fieldValue2)
		}
	}

	return nil
}

func executeCreateUserRequest(response *httptest.ResponseRecorder, request *http.Request) (map[string]interface{}, error) {
	var jsonResponse map[string]interface{}
	CreateUser(response, request)

	responseResult := response.Result()
	defer responseResult.Body.Close()

	if err := json.NewDecoder(responseResult.Body).Decode(&jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}
	return jsonResponse, nil
}

func prepareHTTPRequest(method string, url string, userData database.UserInfo) (*httptest.ResponseRecorder, *http.Request, error) {
	userJSON, err := json.Marshal(userData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal user data to JSON: %v", err)
	}

	request := httptest.NewRequest(method, url, bytes.NewBuffer(userJSON))
	response := httptest.NewRecorder()

	return response, request, nil
}

func getUserDataByEmail(email string) (*database.UserInfo, error) {
	userId, err := validation.GetUserID(database.DB, email, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	userDataFromDB, err := database.SelectUser(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %v", err)
	}
	return userDataFromDB, nil
}
