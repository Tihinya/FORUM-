package session

import (
	"net/http"
	"sync"
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type Session struct {
	UserId     int
	Token      string
	ExpireTime time.Time
	LifeTime   time.Duration
}

type Storage struct {
	storage map[string]Session
	lock    sync.Mutex
}

var SessionStorage Storage

const cookieLifeTime = time.Hour * 24 * 365 // 1 year
const sessionLifeTime = time.Second * 300   // 5 minute

func (s *Storage) CreateSession(userId int) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	sessionToken := uuid.Must(uuid.NewV4()).String()
	s.storage[sessionToken] = Session{
		UserId:     userId,
		Token:      sessionToken,
		ExpireTime: time.Now(),
		LifeTime:   sessionLifeTime,
	}
	return sessionToken
}

func (s *Storage) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
		Path:   "/",
	})
}

func (s *Storage) SetCookie(sessionToken string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(cookieLifeTime),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func ValidateToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	sessionId := cookie.Value
	return sessionId, nil
}

func (s *Storage) GetSession(token string, r *http.Request) Session {
	s.lock.Lock()
	defer s.lock.Unlock()

	session, ok := s.storage[token]
	if !ok {
		return Session{}
	}

	if session.ExpireTime.Before(time.Now()) {
		delete(s.storage, token)
		return Session{}
	}

	return session
}

func (s *Session) RemoveSession() {
	SessionStorage.lock.Lock()
	defer SessionStorage.lock.Unlock()

	delete(SessionStorage.storage, s.Token)
}

func (s *Storage) checkSession() {
	for {
		time.Sleep(time.Second * 60)

		for _, session := range SessionStorage.storage {
			if time.Since(session.ExpireTime) >= session.LifeTime {
				session.RemoveSession()
			}
		}
	}
}

// func (s *Storage) Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("sessionID")
// 		if err != nil || cookie.Value == "" {
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			return
// 		}

// 		sessionID := cookie.Value
// 		session := s.GetSession(sessionID)
// 		if session == nil {
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func init() {
	SessionStorage = Storage{
		storage: make(map[string]Session),
	}

	go SessionStorage.checkSession()
}
