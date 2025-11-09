package router

import (
	"net/http"

	"myapp/config"
	"myapp/internal/app/handler"
	"myapp/internal/app/middleware"
	"myapp/internal/app/session"
)

func NewRouter(u *handler.UserHandler, a *handler.AuthHandler, store session.SessionStore, cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	// Public endpoints
	mux.HandleFunc("/api/register", u.Register)
	mux.HandleFunc("/api/login", a.Login)
	mux.HandleFunc("/api/logout", a.Logout)
	// Protected endpoints
	mux.Handle("/api/users", middleware.WithSession(http.HandlerFunc(u.GetUsers), store, cfg.SessionCookieName))
	return mux
}
