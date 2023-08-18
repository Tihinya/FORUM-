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
	lock    sync.RWMutex
}

var SessionStorage Storage

const cookieLifeTime = time.Hour * 24 * 365 // 1 year
const sessionLifeTime = time.Second * 300   // 5 minutes

func (s *Storage) CreateSession(userId int) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	sessionToken := uuid.Must(uuid.NewV4()).String()
	expireTime := time.Now().Add(sessionLifeTime)
	s.storage[sessionToken] = Session{
		UserId:     userId,
		Token:      sessionToken,
		ExpireTime: expireTime,
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
		HttpOnly: true,
	})
}

func ValidateToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	sessionID := cookie.Value
	return sessionID, nil
}

func (s *Storage) GetSession(token string) Session {
	s.lock.RLock()
	defer s.lock.RUnlock()

	session, ok := s.storage[token]
	if !ok || session.ExpireTime.Before(time.Now()) {
		return Session{}
	}

	return session
}

func (s *Session) RemoveSession() {
	SessionStorage.lock.Lock()
	defer SessionStorage.lock.Unlock()

	delete(SessionStorage.storage, s.Token)
}

func (s *Storage) startSessionCleanupRoutine() {
	for {
		time.Sleep(time.Second * 60)

		s.lock.Lock()
		for token, session := range s.storage {
			if time.Since(session.ExpireTime) >= session.LifeTime {
				delete(s.storage, token)
			}
		}
		s.lock.Unlock()
	}
}

func init() {
	SessionStorage = Storage{
		storage: make(map[string]Session),
	}

	go SessionStorage.startSessionCleanupRoutine()
}