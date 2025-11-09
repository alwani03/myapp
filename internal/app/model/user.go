package model

type User struct {
    ID           int64  `json:"id"`
    Username     string `json:"username"`
    Email        string `json:"email"`
    Name         string `json:"name,omitempty"`
    PasswordHash string `json:"-"`
}