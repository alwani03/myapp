package repository

import (
    "database/sql"
    "errors"
    "strings"

    "myapp/internal/app/model"
)

type SQLUserRepository struct {
    db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) UserRepository {
    return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindAll() []model.User {
    rows, err := r.db.Query(`SELECT id, username, email, COALESCE(NULLIF(name, ''), '') FROM users ORDER BY id`)
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

func (r *SQLUserRepository) Create(u model.User) (model.User, error) {
    u.Username = strings.TrimSpace(u.Username)
    u.Email = strings.TrimSpace(u.Email)
    u.Name = strings.TrimSpace(u.Name)
    if u.Username == "" || u.Email == "" || u.PasswordHash == "" {
        return model.User{}, errors.New("invalid user data")
    }
    var id int64
    err := r.db.QueryRow(
        `INSERT INTO users (username, email, password_hash, name) VALUES ($1, $2, $3, NULLIF($4, '')) RETURNING id`,
        u.Username, u.Email, u.PasswordHash, u.Name,
    ).Scan(&id)
    if err != nil {
        return model.User{}, err
    }
    u.ID = id
    return u, nil
}

func (r *SQLUserRepository) FindByUsername(username string) (model.User, error) {
    var u model.User
    err := r.db.QueryRow(`SELECT id, username, email, COALESCE(name, '') FROM users WHERE username = $1`, username).Scan(&u.ID, &u.Username, &u.Email, &u.Name)
    if err != nil {
        return model.User{}, err
    }
    return u, nil
}

func (r *SQLUserRepository) FindByEmail(email string) (model.User, error) {
    var u model.User
    err := r.db.QueryRow(`SELECT id, username, email, COALESCE(name, '') FROM users WHERE email = $1`, email).Scan(&u.ID, &u.Username, &u.Email, &u.Name)
    if err != nil {
        return model.User{}, err
    }
    return u, nil
}