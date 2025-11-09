package service

import (
	"errors"
	"strings"

	"myapp/internal/app/model"
	"myapp/internal/app/repository"

	"golang.org/x/crypto/bcrypt"
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

type RegisterParams struct {
    Username string
    Email    string
    Password string
    Name     string
}

type LoginParams struct {
    UsernameOrEmail string
    Password        string
}

func (s *UserService) RegisterUser(p RegisterParams) (model.User, error) {
	username := strings.TrimSpace(p.Username)
	email := strings.TrimSpace(p.Email)
	if username == "" || email == "" || strings.TrimSpace(p.Password) == "" {
		return model.User{}, errors.New("username, email, and password are required")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	u := model.User{
		Username:     username,
		Email:        email,
		Name:         strings.TrimSpace(p.Name),
		PasswordHash: string(hash),
	}
	// In current simplified service, we rely on DB unique constraints for duplicates
	created, err := s.repo.Create(u)
	if err != nil {
		return model.User{}, err
	}
	// Do not return password hash in higher layers
	created.PasswordHash = ""
	return created, nil
}

// Authenticate verifies credentials against stored password hash.
// It tries username first; if not found, it tries email.
func (s *UserService) Authenticate(p LoginParams) (model.User, error) {
    ident := strings.TrimSpace(p.UsernameOrEmail)
    pass := strings.TrimSpace(p.Password)
    if ident == "" || pass == "" {
        return model.User{}, errors.New("username/email and password are required")
    }
    var u model.User
    var err error
    // Heuristic: if it looks like an email, try email first; otherwise username first
    if strings.Contains(ident, "@") {
        u, err = s.repo.FindByEmail(ident)
        if err != nil {
            // fallback to username
            u, err = s.repo.FindByUsername(ident)
        }
    } else {
        u, err = s.repo.FindByUsername(ident)
        if err != nil {
            // fallback to email
            u, err = s.repo.FindByEmail(ident)
        }
    }
    if err != nil {
        return model.User{}, errors.New("invalid credentials")
    }
    if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass)) != nil {
        return model.User{}, errors.New("invalid credentials")
    }
    // sanitize
    u.PasswordHash = ""
    return u, nil
}
