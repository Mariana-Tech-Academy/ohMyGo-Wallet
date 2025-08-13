package services

import (
	"errors"
	"math/rand"
	"time"
	"vaqua/models"
	"vaqua/repository"
	"vaqua/utils"
)

type UserService struct {
	Repo repository.UserRepository
}

func (s *UserService) CreateUser(email, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	accountNumber := s.generateAccountNumber()

	user := &models.User{
		Email:         email,
		Password:      hashedPassword,
		AccountNumber: accountNumber,
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	err = utils.ComparePassword(user.Password, password)
	if err != nil {
		return nil, errors.New("invalid Credentials")
	}

	return user, nil
}

func (s *UserService) generateAccountNumber() uint {
	rand.Seed(time.Now().UnixNano())
	return uint(rand.Intn(9000000000) + 1000000000) // 10-digit
}
