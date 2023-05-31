package session

import (
	"net/http"
	"sync"
	"time"

	uuid "github.com/gofrs/uuid"
)

type Session struct {
	UserId      int
	Token       string
	ExpiresTime time.Time
	LifeTime    time.Duration
}

// var sessionStorage = map[string]Session{}

type Storage struct {
	storage map[string]Session
	lock    sync.Mutex
}

var SessionStorage Storage

const cookieLifeTime = time.Hour * 24 * 365 // 1 year
const sessionLifeTime = time.Minute * 5

func (s *Storage) CreateSession(userId int, w http.ResponseWriter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	sessionToken := uuid.Must(uuid.NewV4()).String()
	s.storage[sessionToken] = Session{
		UserId:      userId,
		Token:       sessionToken,
		ExpiresTime: time.Now(),
		LifeTime:    sessionLifeTime,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(cookieLifeTime),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

// func DeleteToken(w http.ResponseWriter) {
// 	http.SetCookie(w, &http.Cookie{
// 		Name:   "session-Id",
// 		MaxAge: -1,
// 		Path:   "/",
// 	})
// }

// func ValidateToken(r *http.Request) (string, error) {
// 	cookie, err := r.Cookie("session_token")
// 	if err != nil {
// 		return "", err
// 	}

// 	sessionId := cookie.Value
// 	return sessionId, nil
// }

// func (s *Storage) CheckSession(token string) {
// 	if time.Since(s.storage[token].ExpiresTime) >= s.storage[token].LifeTime {
// 		delete(SessionStorage.storage, token)
// 	}
// }

func init() {
	SessionStorage = Storage{
		storage: make(map[string]Session),
	}
}
