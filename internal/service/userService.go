package service

import (
	"database/sql"
	"errors"
	"strconv"

	"main.go/internal/helpers"
	"main.go/internal/repository"
)

var ErrRequiredFields = errors.New("required fields are empty")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUsernameTaken = errors.New("username already taken")

type UserService struct {
	Repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) validateCredentials(name, password string) error {
	if name == "" || password == "" {
		return ErrRequiredFields
	}
	return nil
}

func (s *UserService) Login(name, password string) (string, error) {
	if err := s.validateCredentials(name, password); err != nil {
		return "", err
	}
	user, err := s.Repo.GetUserByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	passwordMatch, err := helpers.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !passwordMatch {
		return "", ErrInvalidCredentials
	}
	token, err := helpers.CreateToken(&user)
	if err != nil {
		return "", err
	}
	return token, nil

}
func (s *UserService) Register(name, password string) (string, error) {
	if err := s.validateCredentials(name, password); err != nil {
		return "", err
	}
	exists, err := s.Repo.CheckDuplicateUser(name)
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrUsernameTaken
	}
	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return "", err
	}
	id, err := s.Repo.InsertUser(name, passwordHash)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil

}
