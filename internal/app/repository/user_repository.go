package repository

import "myapp/internal/app/model"

type UserRepository interface {
	FindAll() []model.User
}

type InMemoryUserRepository struct{}

func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) FindAll() []model.User {
	return []model.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
}
