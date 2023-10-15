package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"forum/config"
	ct "forum/controllers"
	"forum/database"
	"forum/login"
	"forum/router"
	"forum/session"
	"forum/validation"
)

var (
	requests = make(map[string]*requestInfo)
	l        sync.Mutex
)

type requestInfo struct {
	count     int
	timestamp time.Time
}

func Auth() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, sessionTokenFound := ct.CheckForSessionToken(r)

			if !sessionTokenFound {
				ct.ReturnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "error")
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
func ModeratorOnly() router.Middleware {
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
			moderatorID, err := database.GetRoleId("moderator")
			if err != nil {
				ct.ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
				return
			}
			if user.RoleID != moderatorID {
				ct.ReturnMessageJSON(w, "Insufficient privileges", http.StatusForbidden, "error")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiter() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now()

			l.Lock()

			info, ok := requests[ip]
			if ok {
				// Check if it's been more than a minute since the first request
				if now.Sub(info.timestamp) > time.Minute {
					info.count = 1
					info.timestamp = now
				} else if info.count >= 200 {
					l.Unlock()
					ct.ReturnMessageJSON(
						w,
						"You have reached the bandwidth limit of your 3G package, your rate has been limited.",
						http.StatusTooManyRequests,
						"error",
					)
					return
				} else {
					info.count++
				}
			} else {
				requests[ip] = &requestInfo{count: 1, timestamp: now}
			}
			l.Unlock() // Unlock after modifying the map

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
		Origin:      "https://localhost:3000",
		Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		Headers:     []string{"Content-Type", "Authorization"},
		Credentials: true,
	}))

	r.AddGlobalMiddleware(RateLimiter())

	// Rate limiter
	r.NewRoute("GET", `/ratelimited`, ct.RateLimited)

	// User
	r.NewRoute("POST", `/user/create`, ct.CreateUser)
	r.NewRoute("GET", `/user/(?P<id>\d+)/get`, ct.ReadUser, AdminOnly())
	r.NewRoute("GET", `/users/get`, ct.ReadUsers)
	r.NewRoute("PATCH", `/user/(?P<id>\d+)/update`, ct.UpdateUser, AdminOnly())
	r.NewRoute("DELETE", `/user/(?P<id>\d+)/delete`, ct.DeleteUser, Auth())
	r.NewRoute("GET", `/user/liked`, ct.ReadUserLikedPosts, Auth())
	r.NewRoute("GET", `/user/disliked`, ct.ReadUserDislikedPosts, Auth())
	r.NewRoute("GET", `/user/likedComments`, ct.ReadUserLikedComments, Auth())
	r.NewRoute("GET", `/user/dislikedComments`, ct.ReadUserDislikedComments, Auth())
	r.NewRoute("GET", `/user/posts`, ct.ReadUserCreatedPosts, Auth())
	r.NewRoute("GET", `/user/createdcomments`, ct.ReadUserCreatedComments, Auth())
	r.NewRoute("GET", `/user/comments`, ct.ReadUserCommentdPosts, Auth())

	// Notifications
	r.NewRoute("GET", `/notifications`, ct.GetNotifications, Auth())
	r.NewRoute("POST", `/readnotifications`, ct.MarkNotificationsRead, Auth())

	// Post
	r.NewRoute("POST", `/post`, ct.CreatePost, Auth())
	r.NewRoute("GET", `/post/(?P<id>\d+)`, ct.ReadPost)
	r.NewRoute("GET", `/posts`, ct.ReadPosts)
	r.NewRoute("PATCH", `/post/(?P<id>\d+)`, ct.UpdatePost, Auth())
	r.NewRoute("DELETE", `/post/(?P<id>\d+)`, ct.DeletePost, Auth())
	r.NewRoute("DELETE", `/post/(?P<id>\d+)/mod`, ct.DeletePostModerator, ModeratorOnly())
	r.NewRoute("GET", `/categories`, ct.ReadCategories)
	r.NewRoute("GET", `/postcategories`, ct.ReadPostCategories)

	// Comment
	r.NewRoute("POST", `/comment/(?P<id>\d+)`, ct.CreateComment, Auth())
	r.NewRoute("GET", `/comment/(?P<id>\d+)`, ct.ReadComment)
	r.NewRoute("GET", `/comments/(?P<id>\d+)`, ct.ReadComments)
	r.NewRoute("PATCH", `/comment/(?P<id>\d+)`, ct.UpdateComment, Auth())
	r.NewRoute("DELETE", `/comment/(?P<id>\d+)`, ct.DeleteComment, Auth())

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
	r.NewRoute("POST", `/logout`, ct.LogOut)
	r.NewRoute("GET", `/login/google`, ct.GoogleLogin)
	r.NewRoute("GET", `/login/google/callback`, ct.GoogleCallback)
	r.NewRoute("GET", `/login/github`, ct.GithubLogin)
	r.NewRoute("GET", `/login/github/callback`, ct.GithubCallback)
	r.NewRoute("GET", `/authorized`, ct.CheckAuthorization, Auth())

	r.NewRoute("GET", `/login/github/redirect`, ct.GithubCallbackRedirect)

	//Role
	r.NewRoute("POST", `/promotion/(?P<id>\d+)`, ct.Promotion)
	r.NewRoute("GET", `/promotions/get`, ct.ReadPromotions, AdminOnly())
	r.NewRoute("PATCH", `/promotion/(?P<id>\d+)/update`, ct.UpdatePromotion, AdminOnly())
	r.NewRoute("DELETE", `/promotion/(?P<id>\d+)/delete`, ct.DeletePromotion, AdminOnly())

	log.Println("Ctrl + Click on the link: https://localhost:" + config.Config.Port)
	log.Println("To stop the server press `Ctrl + C`")

	log.Fatal(http.ListenAndServeTLS(":"+config.Config.Port, "cert.pem", "key.pem", nil))
}
