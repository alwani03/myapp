package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"myapp/config"
	"myapp/internal/app/session"
)

type AuthHandler struct {
	store session.SessionStore
	cfg   config.Config
}

func NewAuthHandler(store session.SessionStore, cfg config.Config) *AuthHandler {
	return &AuthHandler{store: store, cfg: cfg}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type messageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Username != h.cfg.AdminUser || req.Password != h.cfg.AdminPassword {
		http.Error(w, "gagal login", http.StatusUnauthorized)
		return
	}

	sess, err := h.store.Create(req.Username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.SessionCookieName,
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  sess.ExpiresAt,
		// Secure: true, // enable if serving over HTTPS
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(messageResponse{Message: "login successful", Success: true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	c, err := r.Cookie(h.cfg.SessionCookieName)
	if err == nil && c.Value != "" {
		h.store.Delete(c.Value)
	}
	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(messageResponse{Message: "logout successful", Success: true})
}
