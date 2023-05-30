package login

import (
	"forum/session"
	"net/http"
)

var saveSession = make(map[int]string)

func AddLogin(w http.ResponseWriter, userId int) {
	session.SessionStorage.CreateSession(userId, w)
}
func Registration(w http.ResponseWriter, r *http.Request) {

}
