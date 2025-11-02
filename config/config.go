package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
    Port       int
    JwtSecret  string
    AdminUser  string
    AdminPassword string
    SessionCookieName string
    SessionTTLMinutes int
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    DBDriver   string
}

func Load() Config {
	// Load .env file and override existing env values to ensure latest changes are used
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	port := 55000
	if v := os.Getenv("APP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	} else if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "changeme"
    }

    adminUser := os.Getenv("ADMIN_USER")
    if adminUser == "" {
        adminUser = "admin"
    }
    adminPass := os.Getenv("ADMIN_PASSWORD")
    if adminPass == "" {
        adminPass = "secret"
    }
    sessionCookie := os.Getenv("SESSION_COOKIE_NAME")
    if sessionCookie == "" {
        sessionCookie = "session_id"
    }
    sessionTTL := 60
    if v := os.Getenv("SESSION_TTL_MINUTES"); v != "" {
        if p, err := strconv.Atoi(v); err == nil {
            sessionTTL = p
        }
    }
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	dbPort := 5432
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			dbPort = p
		}
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql" // default
	}

    log.Printf("Loaded config: Port=%d, JWT_SECRET=%s, Admin=%s, SessionCookie=%s, SessionTTL=%d, DB=%s:%d/%s, DB_DRIVER=%s", port, secret, adminUser, sessionCookie, sessionTTL, host, dbPort, name, driver)
    return Config{
        Port:       port,
        JwtSecret:  secret,
        AdminUser:  adminUser,
        AdminPassword: adminPass,
        SessionCookieName: sessionCookie,
        SessionTTLMinutes: sessionTTL,
        DBHost:     host,
        DBPort:     dbPort,
        DBUser:     user,
        DBPassword: pass,
        DBName:     name,
        DBDriver:   driver,
    }
}
