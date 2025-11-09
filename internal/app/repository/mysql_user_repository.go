package repository

import (
    "database/sql"
    "errors"
    "strings"

    "myapp/internal/app/model"
)

type MySQLUserRepository struct {
    db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) UserRepository {
    return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) FindAll() []model.User {
    rows, err := r.db.Query(`SELECT id, username, email, IFNULL(name, '') FROM users ORDER BY id`)
    if err != nil {
        return []model.User{}
    }
    defer rows.Close()
    users := make([]model.User, 0)
    for rows.Next() {
        var u model.User
        if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Name); err == nil {
            users = append(users, u)
        }
    }
    _ = rows.Err()
    return users
}

func (r *MySQLUserRepository) Create(u model.User) (model.User, error) {
    u.Username = strings.TrimSpace(u.Username)
    u.Email = strings.TrimSpace(u.Email)
    u.Name = strings.TrimSpace(u.Name)
    if u.Username == "" || u.Email == "" || u.PasswordHash == "" {
        return model.User{}, errors.New("invalid user data")
    }
    res, err := r.db.Exec(
        `INSERT INTO users (username, email, password_hash, name) VALUES (?, ?, ?, NULLIF(?, ''))`,
        u.Username, u.Email, u.PasswordHash, u.Name,
    )
    if err != nil {
        return model.User{}, err
    }
    id, _ := res.LastInsertId()
    u.ID = id
    return u, nil
}

func (r *MySQLUserRepository) FindByUsername(username string) (model.User, error) {
    var u model.User
    err := r.db.QueryRow(`SELECT id, username, email, IFNULL(name, ''), password_hash FROM users WHERE username = ?`, username).Scan(&u.ID, &u.Username, &u.Email, &u.Name, &u.PasswordHash)
    if err != nil {
        return model.User{}, err
    }
    return u, nil
}

func (r *MySQLUserRepository) FindByEmail(email string) (model.User, error) {
    var u model.User
    err := r.db.QueryRow(`SELECT id, username, email, IFNULL(name, ''), password_hash FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Username, &u.Email, &u.Name, &u.PasswordHash)
    if err != nil {
        return model.User{}, err
    }
    return u, nil
}