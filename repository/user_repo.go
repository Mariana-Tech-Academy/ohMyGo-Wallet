package repository

import (
	"vaqua/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
