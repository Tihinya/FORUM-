package login

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

var saveSession = make(map[int]string)

func AddLogin(w http.ResponseWriter, userId int) error {
	//Create token
	UUIDtoken, err := uuid.NewV4()
	if err != nil {
		return err
	}
	sessionToken := UUIDtoken.String()
	//set up cookies for web
	oneYear := time.Now().Add(time.Hour * 8760)
	cookie := &http.Cookie{
		Name:    "session-Id",
		Value:   sessionToken,
		Expires: oneYear,
	}
	//save session to map
	saveSession[userId] = sessionToken
	// git cookies to user
	http.SetCookie(w, cookie)

	return nil
}
func Registration(w http.ResponseWriter, r *http.Request) {

}
