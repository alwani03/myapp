package middleware

import (
    "net/http"
    "strings"

    "myapp/internal/app/session"
)

// WithSession protects a handler by requiring a valid session cookie.
// Routes like /api/login should not be wrapped with this.
func WithSession(next http.Handler, store session.SessionStore, cookieName string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow CORS preflight and public routes pattern if needed
        if r.Method == http.MethodOptions {
            next.ServeHTTP(w, r)
            return
        }
        // Read session cookie
        c, err := r.Cookie(cookieName)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        sid := strings.TrimSpace(c.Value)
        if sid == "" || !store.IsValid(sid) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}