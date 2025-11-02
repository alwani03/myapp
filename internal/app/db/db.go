package db

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq"
    "myapp/config"
)

func ConnectPG(cfg config.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    db.SetConnMaxLifetime(0)
    db.SetMaxIdleConns(5)
    db.SetMaxOpenConns(10)

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}