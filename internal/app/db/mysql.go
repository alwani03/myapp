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
	// DSN format: user:pass@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, err := sql.Open("mysql", dsn)
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
