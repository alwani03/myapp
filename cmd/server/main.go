package main

import (
	"database/sql"
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
	var dbConn *sql.DB
	var repo repository.UserRepository
	switch cfg.DBDriver {
	case "mysql":
		if db, err := appdb.ConnectMySQL(cfg); err != nil {
			log.Printf("DB (mysql) connection failed: %v (using in-memory repo)", err)
		} else {
			log.Printf("DB (mysql) connected successfully")
			dbConn = db
			if err := appdb.RunMigrations(dbConn, "mysql"); err != nil {
				log.Printf("DB migrations failed: %v", err)
			}
			defer db.Close()
		}
	default:
		if db, err := appdb.ConnectPG(cfg); err != nil {
			log.Printf("DB (postgres) connection failed: %v (using in-memory repo)", err)
		} else {
			log.Printf("DB (postgres) connected successfully")
			dbConn = db
			if err := appdb.RunMigrations(dbConn, "postgres"); err != nil {
				log.Printf("DB migrations failed: %v", err)
			}
			defer db.Close()
		}
	}
	if dbConn != nil {
		if cfg.DBDriver == "mysql" {
			repo = repository.NewMySQLUserRepository(dbConn)
		} else {
			repo = repository.NewSQLUserRepository(dbConn)
		}
	} else {
		repo = repository.NewUserRepository()
	}
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)
	store := session.NewInMemorySessionStore(cfg.SessionTTLMinutes)
    auth := handler.NewAuthHandler(store, cfg, svc)

	mux := router.NewRouter(h, auth, store, cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
