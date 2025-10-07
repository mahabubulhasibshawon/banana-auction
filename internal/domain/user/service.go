package user

import (
	"banana-auction/internal/infrastructure/utils"
	"errors"
)

type Service interface {
	Register(username, password, name, role string) (int, error)
	Login(username, password string) (string, error)
	GetUser(id int) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(username, password, name, role string) (int, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	u := User{
		Username:     username,
		PasswordHash: hashedPassword,
		Name:         name,
		Role:         role,
	}

	return s.repo.Create(u)
}

func (s *service) Login(username, password string) (string, error) {
	u, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, u.PasswordHash) {
		return "", errors.New("invalid password")
	}

	return utils.GenerateJWT(u.ID)
}

func (s *service) GetUser(id int) (User, error) {
	return s.repo.GetByID(id)
}
