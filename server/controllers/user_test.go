package controllers

import (
	"forum/database"
	"os"
	"testing"
)

const testDBPath = "mock_database.db"

var validUserTestData = database.UserInfo{
	Username:             "testuser",
	Email:                "test1user@example.com",
	Password:             "testpassword",
	PasswordConfirmation: "testpassword",
	Gender:               "male",
	Age:                  "12",
}

var notValidUserTestData = database.UserInfo{
	Username:             "testuser",
	Email:                "test2user@example.com",
	Password:             "testpassword",
	PasswordConfirmation: "testpassword",
	Gender:               "male",
}

var ignoreList = map[string]bool{"ID": true, "Password": true, "PasswordConfirmation": true}

func TestValidUserRegistration(t *testing.T) {
	if err := database.OpenDatabase(testDBPath); err != nil {
		t.Fatalf("Failed opening db: %v", err)
	}
	defer cleanupDB(t)

	response, request, err := prepareHTTPRequest("POST", "/create-user", validUserTestData)
	if err != nil {
		t.Fatal(err)
	}
	jsonResponse, err := executeCreateUserRequest(response, request)
	if err != nil {
		t.Fatal(err)
	}

	if jsonResponse["status"] != "success" {
		t.Fatalf("Status %s. %s", jsonResponse["status"], jsonResponse["message"])
	}

	userDataFromDB, err := getUserDataByEmail(validUserTestData.Email)
	if err != nil {
		t.Fatal(err)
	}
	err = compareTwoStructs(*userDataFromDB, validUserTestData, ignoreList)
	if err != nil {
		t.Fatal(err)
	}

}

func TestInvalidUserRegistration(t *testing.T) {
	if err := database.OpenDatabase(testDBPath); err != nil {
		t.Fatalf("Failed opening db: %v", err)
	}
	defer cleanupDB(t)

	response, request, err := prepareHTTPRequest("POST", "/create-user", notValidUserTestData)
	if err != nil {
		t.Fatal(err)
	}
	jsonResponse, err := executeCreateUserRequest(response, request)
	if err != nil {
		t.Fatal(err)
	}

	if jsonResponse["status"] != "success" {
		expectedErrMsg := "Invalid age, please select age"
		if jsonResponse["message"] != expectedErrMsg {
			t.Fatalf("Expected error message: %s, got: Status %s. %s", expectedErrMsg, jsonResponse["status"], jsonResponse["message"])
		}
	}

}

func TestRandomUserRegistration(t *testing.T) {
	if err := database.OpenDatabase(testDBPath); err != nil {
		t.Fatalf("Failed opening db: %v", err)
	}
	defer cleanupDB(t)

	userAmount := 300
	for userNumber := 0; userNumber <= userAmount; userNumber++ {
		userData := randomUserData(userNumber)

		response, request, err := prepareHTTPRequest("POST", "/create-user", userData)
		if err != nil {
			t.Fatal(err)
		}
		jsonResponse, err := executeCreateUserRequest(response, request)
		if err != nil {
			t.Fatal(err)
		}

		if jsonResponse["status"] != "success" {
			t.Fatalf("Status %s. %s", jsonResponse["status"], jsonResponse["message"])
		}

		userDataFromDB, err := getUserDataByEmail(userData.Email)
		if err != nil {
			t.Fatal(err)
		}
		err = compareTwoStructs(*userDataFromDB, userData, ignoreList)
		if err != nil {
			t.Fatal(err)
		}

	}
}

func cleanupDB(t *testing.T) {
	err := os.Remove(testDBPath)
	if err != nil {
		t.Fatalf("Failed to remove database: %v", err)
	}
}
