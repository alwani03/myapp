package service

import (
    "myapp/internal/app/model"
    "myapp/internal/app/repository"
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) ListUsers() []model.User {
    return s.repo.FindAll()
}