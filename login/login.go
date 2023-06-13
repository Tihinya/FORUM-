package login

import (
	"forum/session"
	"net/http"
)

func AddLogin(w http.ResponseWriter, userId int) {
	token := session.SessionStorage.CreateSession(userId)
	session.SessionStorage.SetCookie(token, w)
}

func Registration(w http.ResponseWriter, r *http.Request) {

}
