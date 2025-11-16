package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"myapp/config"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL(cfg config.Config) (*sql.DB, error) {
    base := "parseTime=true&charset=utf8mb4&loc=Local&multiStatements=true"
    if cfg.DBOptions != "" {
        base = cfg.DBOptions + "&" + base
    }
    host := cfg.DBHost
    if host == "localhost" {
        host = "127.0.0.1"
    }
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
        cfg.DBUser, cfg.DBPassword, host, cfg.DBPort, cfg.DBName, base,
    )
    db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
    db.SetConnMaxLifetime(0)
    db.SetMaxIdleConns(5)
    db.SetMaxOpenConns(10)
    for i := 0; i < 10; i++ {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        err = db.PingContext(ctx)
        cancel()
        if err == nil {
            return db, nil
        }
        time.Sleep(1 * time.Second)
    }
    db.Close()
    if err != nil {
        return nil, err
    }
    return db, nil
}
