package main

import (
	"log"
	"net/http"

	"forum/config"
	ct "forum/controllers"
	"forum/database"
	"forum/router"
)

func Auth() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionToken, sessionTokenFound := ct.CheckForSessionToken(r)
			if !sessionTokenFound {
				ct.ReturnMessageJSON(w, "Only authorized users can give likes", http.StatusUnauthorized, "error")
				log.Println("Auth middleware fail")
				return
			}

			if !ct.CheckIfUserLoggedin(sessionToken) {
				ct.ReturnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
				log.Println("Auth middleware fail")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	r := router.NewRouter()

	database.CreateTables()

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
	r.NewRoute("GET", `/user/liked`, ct.ReadUserLikedPosts, Auth())
	r.NewRoute("GET", `/user/disliked`, ct.ReadUserDislikedPosts, Auth())
	r.NewRoute("GET", `/user/likedComments`, ct.ReadUserLikedComments, Auth())
	r.NewRoute("GET", `/user/dislikedComments`, ct.ReadUserDislikedComments, Auth())
	r.NewRoute("GET", `/user/posts`, ct.ReadUserCreatedPosts, Auth())
	r.NewRoute("GET", `/user/comments`, ct.ReadUserCommentdPosts)

	// Post
	r.NewRoute("POST", `/post`, ct.CreatePost, Auth())
	r.NewRoute("GET", `/post/(?P<id>\d+)`, ct.ReadPost)
	r.NewRoute("PATCH", `/post/(?P<id>\d+)`, ct.UpdatePost, Auth())
	r.NewRoute("DELETE", `/post/(?P<id>\d+)`, ct.DeletePost, Auth())
	r.NewRoute("GET", `/posts`, ct.ReadPosts)
	r.NewRoute("GET", `/categories`, ct.ReadCategories)
	r.NewRoute("GET", `/postcategories`, ct.ReadPostCategories)

	// Comment
	r.NewRoute("POST", `/comment/(?P<id>\d+)`, ct.CreateComment, Auth())
	r.NewRoute("GET", `/comment/(?P<id>\d+)`, ct.ReadComment)
	r.NewRoute("PATCH", `/comment/(?P<id>\d+)`, ct.UpdateComment, Auth())
	r.NewRoute("DELETE", `/comment/(?P<id>\d+)`, ct.DeleteComment, Auth())
	r.NewRoute("GET", `/comments/(?P<id>\d+)`, ct.ReadComments)

	// Like
	r.NewRoute("POST", `/post/(?P<id>\d+)/like`, ct.LikePost, Auth())
	r.NewRoute("POST", `/post/(?P<id>\d+)/unlike`, ct.UnlikePost, Auth())
	r.NewRoute("POST", `/comment/(?P<id>\d+)/like`, ct.LikeComment, Auth())
	r.NewRoute("POST", `/comment/(?P<id>\d+)/unlike`, ct.UnlikeComment, Auth())

	// Dislike
	r.NewRoute("POST", `/post/(?P<id>\d+)/dislike`, ct.DislikePost, Auth())
	r.NewRoute("POST", `/post/(?P<id>\d+)/undislike`, ct.UndislikePost, Auth())
	r.NewRoute("POST", `/comment/(?P<id>\d+)/dislike`, ct.DislikeComment, Auth())
	r.NewRoute("POST", `/comment/(?P<id>\d+)/undislike`, ct.UndislikeComment, Auth())

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
