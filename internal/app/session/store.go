package session

import (
    "crypto/rand"
    "encoding/hex"
    "sync"
    "time"
)

type Session struct {
    ID        string
    User      string
    ExpiresAt time.Time
}

type SessionStore interface {
    Create(user string) (Session, error)
    Get(id string) (Session, bool)
    Delete(id string)
    IsValid(id string) bool
}

type InMemorySessionStore struct {
    mu         sync.RWMutex
    sessions   map[string]Session
    ttlMinutes int
}

func NewInMemorySessionStore(ttlMinutes int) *InMemorySessionStore {
    return &InMemorySessionStore{
        sessions:   make(map[string]Session),
        ttlMinutes: ttlMinutes,
    }
}

func (s *InMemorySessionStore) Create(user string) (Session, error) {
    id, err := randomID(32)
    if err != nil {
        return Session{}, err
    }
    sess := Session{
        ID:        id,
        User:      user,
        ExpiresAt: time.Now().Add(time.Duration(s.ttlMinutes) * time.Minute),
    }
    s.mu.Lock()
    s.sessions[id] = sess
    s.mu.Unlock()
    return sess, nil
}

func (s *InMemorySessionStore) Get(id string) (Session, bool) {
    s.mu.RLock()
    sess, ok := s.sessions[id]
    s.mu.RUnlock()
    if !ok {
        return Session{}, false
    }
    if time.Now().After(sess.ExpiresAt) {
        // expired, cleanup
        s.Delete(id)
        return Session{}, false
    }
    return sess, true
}

func (s *InMemorySessionStore) Delete(id string) {
    s.mu.Lock()
    delete(s.sessions, id)
    s.mu.Unlock()
}

func (s *InMemorySessionStore) IsValid(id string) bool {
    _, ok := s.Get(id)
    return ok
}

func randomID(n int) (string, error) {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}