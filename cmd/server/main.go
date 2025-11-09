package main

import (
	"fmt"
	"log"
	"net/http"

	"myapp/config"
	appdb "myapp/internal/app/db"
	"myapp/internal/app/handler"
	"myapp/internal/app/repository"
	"myapp/internal/app/router"
	"myapp/internal/app/service"
	"myapp/internal/app/session"
)

func main() {
	cfg := config.Load()
	switch cfg.DBDriver {
	case "mysql":
		if db, err := appdb.ConnectMySQL(cfg); err != nil {
			log.Printf("DB (mysql) connection failed: %v (using in-memory repo)", err)
		} else {
			log.Printf("DB (mysql) connected successfully")
			defer db.Close()
		}
	default:
		if db, err := appdb.ConnectPG(cfg); err != nil {
			log.Printf("DB (postgres) connection failed: %v (using in-memory repo)", err)
		} else {
			log.Printf("DB (postgres) connected successfully")
			defer db.Close()
		}
	}
	repo := repository.NewUserRepository()
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)
	store := session.NewInMemorySessionStore(cfg.SessionTTLMinutes)
	auth := handler.NewAuthHandler(store, cfg)

	mux := router.NewRouter(h, auth, store, cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
