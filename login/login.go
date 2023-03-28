package login

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

var saveSession = make(map[int]string)

func Login(w http.ResponseWriter, r *http.Request) {
	login(w, 12)
	fmt.Fprintf(w, "skdfj")
}

func login(w http.ResponseWriter, userId int) {

	//Create token
	UUIDtoken, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	sessionToken := UUIDtoken.String()
	//set up cookies for web
	cookie := &http.Cookie{
		Name:    "session-Id",
		Value:   sessionToken,
		Expires: time.Now().Add(5 * time.Minute),
	}
	//save session to map
	saveSession[userId] = sessionToken
	// git cookies to user
	http.SetCookie(w, cookie)

}
func Registration(w http.ResponseWriter, r *http.Request) {

}
