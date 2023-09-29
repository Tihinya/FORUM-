package main

import (
	"log"
	"net/http"

	"forum/config"
	ct "forum/controllers"
	"forum/database"
	"forum/login"
	"forum/router"
	"forum/session"
	"forum/validation"
)

func Auth() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, sessionTokenFound := ct.CheckForSessionToken(r)
			if !sessionTokenFound {
				ct.ReturnMessageJSON(w, "Only authorized users can give likes", http.StatusUnauthorized, "error")
				log.Println("Auth middleware fail")
				return
			}
			if !ct.CheckIfUserLoggedin(r) {
				ct.ReturnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "error")
				log.Println("Auth middleware fail")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AdminOnly() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the user's session
			session, err := session.SessionStorage.GetSession(r)
			if err != nil {
				ct.ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
				return
			}
			// Check if a valid session exists
			if session == nil {
				ct.ReturnMessageJSON(w, "Unauthorized", http.StatusUnauthorized, "error")
				return
			}
			// Retrieve the user from the database based on the ID
			user, err := database.SelectUser(session.UserId)
			if err != nil {
				ct.ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
				return
			}
			// Check if the user has the admin role
			adminID, err := validation.GetUserID(database.DB, "", "admin")
			if err != nil {
				ct.ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
				return
			}
			if user.ID != adminID {
				ct.ReturnMessageJSON(w, "Insufficient privileges", http.StatusForbidden, "error")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	r := router.NewRouter()

	database.CreateTables()

	roles := []string{"user", "moderator", "admin"}
	for _, role := range roles {
		exist, err := validation.ValidateRole(database.DB, role)
		if err != nil {
			log.Println(err)
			return
		}
		if !exist {
			database.AddRole(role)
		}
	}

	login.CreateAdminUser()

	http.HandleFunc("/", r.ServeWithCORS(router.CORS{
		Origin:      "http://localhost:3000",
		Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Headers:     []string{"Content-Type", "Authorization"},
		Credentials: true,
	}))

	// User
	r.NewRoute("POST", `/user/create`, ct.CreateUser)
	r.NewRoute("GET", `/user/(?P<id>\d+)/get`, ct.ReadUser, AdminOnly())
	r.NewRoute("GET", `/users/get`, ct.ReadUsers)
	r.NewRoute("PATCH", `/user/(?P<id>\d+)/update`, ct.UpdateUser)
	r.NewRoute("DELETE", `/user/(?P<id>\d+)/delete`, ct.DeleteUser)
	r.NewRoute("GET", `/user/liked`, ct.ReadUserLikedPosts)
	r.NewRoute("GET", `/user/disliked`, ct.ReadUserDislikedPosts)
	r.NewRoute("GET", `/user/likedComments`, ct.ReadUserLikedComments)
	r.NewRoute("GET", `/user/dislikedComments`, ct.ReadUserDislikedComments)
	r.NewRoute("GET", `/user/posts`, ct.ReadUserCreatedPosts)
	r.NewRoute("GET", `/user/comments`, ct.ReadUserCommentdPosts)

	// Post
	r.NewRoute("POST", `/post`, ct.CreatePost)
	r.NewRoute("GET", `/post/(?P<id>\d+)`, ct.ReadPost)
	r.NewRoute("PATCH", `/post/(?P<id>\d+)`, ct.UpdatePost)
	r.NewRoute("DELETE", `/post/(?P<id>\d+)`, ct.DeletePost)
	r.NewRoute("GET", `/posts`, ct.ReadPosts)
	r.NewRoute("GET", `/categories`, ct.ReadCategories)
	r.NewRoute("GET", `/postcategories`, ct.ReadPostCategories)

	// Comment
	r.NewRoute("POST", `/comment/(?P<id>\d+)`, ct.CreateComment)
	r.NewRoute("GET", `/comment/(?P<id>\d+)`, ct.ReadComment)
	r.NewRoute("PATCH", `/comment/(?P<id>\d+)`, ct.UpdateComment)
	r.NewRoute("DELETE", `/comment/(?P<id>\d+)`, ct.DeleteComment)
	r.NewRoute("GET", `/comments/(?P<id>\d+)`, ct.ReadComments)

	// Like
	r.NewRoute("POST", `/post/(?P<id>\d+)/like`, ct.LikePost)
	r.NewRoute("POST", `/post/(?P<id>\d+)/unlike`, ct.UnlikePost)
	r.NewRoute("POST", `/comment/(?P<id>\d+)/like`, ct.LikeComment)
	r.NewRoute("POST", `/comment/(?P<id>\d+)/unlike`, ct.UnlikeComment)

	// Dislike
	r.NewRoute("POST", `/post/(?P<id>\d+)/dislike`, ct.DislikePost)
	r.NewRoute("POST", `/post/(?P<id>\d+)/undislike`, ct.UndislikePost)
	r.NewRoute("POST", `/comment/(?P<id>\d+)/dislike`, ct.DislikeComment)
	r.NewRoute("POST", `/comment/(?P<id>\d+)/undislike`, ct.UndislikeComment)

	// Temp
	r.NewRoute("GET", `/likes`, ct.Temp_getLikes)
	r.NewRoute("GET", `/dislikes`, ct.Temp_getDislikes)

	// Login
	r.NewRoute("POST", `/login`, ct.Login)
	r.NewRoute("GET", `/logout`, ct.LogOut)
	r.NewRoute("GET", `/login/google`, ct.GoogleLogin)
	r.NewRoute("GET", `/login/google/callback`, ct.GoogleCallback)
	r.NewRoute("GET", `/login/github`, ct.GithubLogin)
	r.NewRoute("GET", `/login/github/callback`, ct.GithubCallback)
	r.NewRoute("GET", `/login/github/redirect`, ct.GithubCallbackRedirect)

	log.Println("Ctrl + Click on the link: https://localhost:" + config.Config.Port)
	log.Println("To stop the server press `Ctrl + C`")

	http.ListenAndServe(":"+config.Config.Port, nil)
}
