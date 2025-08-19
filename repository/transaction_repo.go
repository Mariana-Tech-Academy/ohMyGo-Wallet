package repository

import (
	"vaqua/db"
	"vaqua/models"
)


type TransactionRepository interface {
	CreateTransaction(tx *models.Transaction) error
	GetAllTransactionsByUser(userID uint, limit int, offset int) ([]models.Transaction, error)
}

type TransactionRepo struct{}

func (r *TransactionRepo) CreateTransaction(tx *models.Transaction) error {
	return db.Db.Create(tx).Error
}


func (r *TransactionRepo) GetAllTransactionsByUser(userID uint, limit int, offset int) ([]models.Transaction, error) {
	var txs []models.Transaction
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if err := db.Db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
