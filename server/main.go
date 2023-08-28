package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/config"
	ct "forum/controllers"
	"forum/database"
	"forum/router"
)

func ExampleMiddleware() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Your middleware logic here
			fmt.Println("Example middleware executed 1")

			// Call the next middleware/handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	r := router.NewRouter()

	database.CreateTables()

	// middleware usage example
	r.AddGlobalMiddleware(ExampleMiddleware())

	http.HandleFunc("/", r.ServeWithCORS(router.CORS{
		Origin:      "http://localhost:3000",
		Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Headers:     []string{"Content-Type", "Authorization"},
		Credentials: true,
	}))

	// User
	r.NewRoute("POST", `/user/create`, ct.CreateUser)
	r.NewRoute("GET", `/user/(?P<id>\d+)/get`, ct.ReadUser)
	r.NewRoute("GET", `/users/get`, ct.ReadUsers)
	r.NewRoute("PATCH", `/user/(?P<id>\d+)/update`, ct.UpdateUser)
	r.NewRoute("DELETE", `/user/(?P<id>\d+)/delete`, ct.DeleteUser)
	r.NewRoute("GET", `/user/liked`, ct.ReadUserLikedPosts)
	r.NewRoute("GET", `/user/disliked`, ct.ReadUserDislikedPosts)
	r.NewRoute("GET", `/user/likedComments`, ct.ReadUserLikedComments)
	r.NewRoute("GET", `/user/dislikedComments`, ct.ReadUserDislikedComments)
	r.NewRoute("GET", `/user/posts`, ct.ReadUserCreatedPosts)

	// Post
	r.NewRoute("POST", `/post`, ct.CreatePost)
	r.NewRoute("GET", `/post/(?P<id>\d+)`, ct.ReadPost)
	r.NewRoute("PATCH", `/post/(?P<id>\d+)`, ct.UpdatePost)
	r.NewRoute("DELETE", `/post/(?P<id>\d+)`, ct.DeletePost)
	r.NewRoute("GET", `/posts`, ct.ReadPosts)
	r.NewRoute("GET", `/categories`, ct.ReadCategories)

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
