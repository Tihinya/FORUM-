package controllers

import (
	"forum/login"
	"forum/session"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	userId := 228 // test user id
	login.AddLogin(w, userId)
}
func LogOut(w http.ResponseWriter, r *http.Request) {
	token, err := session.ValidateToken(r)

	if err != nil {
		return
	}
	s := session.SessionStorage.GetSession(token, r)
	s.RemoveSession()
	session.SessionStorage.DeleteCookie(w)
}
