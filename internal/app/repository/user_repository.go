package repository

import "myapp/internal/app/model"

type UserRepository interface {
    FindAll() []model.User
    Create(u model.User) (model.User, error)
    FindByUsername(username string) (model.User, error)
    FindByEmail(email string) (model.User, error)
}

type InMemoryUserRepository struct{}

func NewUserRepository() UserRepository {
    return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) FindAll() []model.User {
    return []model.User{
        {ID: 1, Username: "alice", Email: "alice@example.com", Name: "Alice"},
        {ID: 2, Username: "bob", Email: "bob@example.com", Name: "Bob"},
    }
}

func (r *InMemoryUserRepository) Create(u model.User) (model.User, error) {
    // simulate ID assignment
    u.ID = 100
    return u, nil
}

func (r *InMemoryUserRepository) FindByUsername(username string) (model.User, error) {
    return model.User{}, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (model.User, error) {
    return model.User{}, nil
}

// SQL-backed implementation provided in sql_user_repository.go
