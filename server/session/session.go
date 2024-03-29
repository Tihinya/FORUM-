package session

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type Session struct {
	UserId     int
	Token      string
	ExpireTime time.Time
}

type Storage struct {
	storage map[string]*Session
	lock    sync.RWMutex
}

var SessionStorage Storage

const (
	cookieLifeTime  = time.Hour * 24 * 365 // 1 year
	sessionLifeTime = time.Second * 300    // 5 minutes
)

// Initialize storage (for testing)
func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]*Session),
		lock:    sync.RWMutex{},
	}
}

func (s *Storage) CreateSession(userId int) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	sessionToken := uuid.Must(uuid.NewV4()).String()
	expireTime := time.Now().Add(sessionLifeTime)
	s.storage[sessionToken] = &Session{
		UserId:     userId,
		Token:      sessionToken,
		ExpireTime: expireTime,
	}
	return sessionToken
}

func (s *Storage) DeleteCookie(w http.ResponseWriter) {
	var cookie http.Cookie = http.Cookie{
		Name:     "session_token",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	fmt.Println(cookie.Valid())

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func (s *Session) GetUID() int {
	return s.UserId
}

func (s *Storage) SetCookie(sessionToken string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(cookieLifeTime),
		MaxAge:   int(cookieLifeTime.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func (s *Storage) validateToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	sessionID := cookie.Value
	return sessionID, nil
}

func (s *Storage) GetSession(r *http.Request) (*Session, error) {
	sessionID, err := s.validateToken(r)
	if err != nil {
		return nil, err
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	session, ok := s.storage[sessionID]
	if !ok || session.ExpireTime.Before(time.Now()) {
		delete(s.storage, sessionID)
		return nil, fmt.Errorf("session does not exist or expired")
	}

	session.ExpireTime = time.Now().Add(sessionLifeTime)
	return session, nil
}

func (s *Session) RemoveSession() {
	SessionStorage.lock.Lock()
	defer SessionStorage.lock.Unlock()

	delete(SessionStorage.storage, s.Token)
}

func (s *Storage) startSessionCleanupRoutine() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.lock.Lock()
		for token, session := range s.storage {
			if session.ExpireTime.Before(time.Now()) {
				delete(s.storage, token)
			}
		}
		s.lock.Unlock()
	}
}

func GetUserId(r *http.Request) (int, error) {
	SessionData, err := SessionStorage.GetSession(r)
	if err != nil {
		return 0, err
	}

	userID := SessionData.UserId

	return userID, err
}

func init() {
	SessionStorage = Storage{
		storage: make(map[string]*Session),
	}

	go SessionStorage.startSessionCleanupRoutine()
}
