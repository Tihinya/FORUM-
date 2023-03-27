package login

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func login(w http.ResponseWriter, r *http.Request) {
	userId := 18

	saveSession := make(map[int]string)
	//Create token
	UUIDtoken, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	sessionToken := UUIDtoken.String()

	cookie := &http.Cookie{
		Name:    "session-Id",
		Value:   sessionToken,
		Expires: time.Now().Add(5 * time.Minute),
	}
	saveSession[userId] = sessionToken

	http.SetCookie(w, cookie)

	fmt.Println(saveSession)

}
