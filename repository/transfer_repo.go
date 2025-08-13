package repository

import (
	"gorm.io/gorm"
	"vaqua/models"
)

type TransferRequestRepository interface {
	Create(request *models.TransferRequest) error
	FindByID(id uint) (*models.TransferRequest, error)
}

type TransferRequestRepo struct {
	DB *gorm.DB
}

func (r *TransferRequestRepo) Create(request *models.TransferRequest) error {
	return r.DB.Create(request).Error
}

func (r *TransferRequestRepo) FindByID(id uint) (*models.TransferRequest, error) {
	var req models.TransferRequest
	err := r.DB.First(&req, id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}
