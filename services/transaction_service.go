package services

import (
	"vaqua/models"
	"vaqua/repository"
)

type TransactionService struct {
	Repo repository.TransactionRepository
}

func (s *TransactionService) CreateTransaction(tx *models.Transaction) error {
	return s.Repo.CreateTransaction(tx)
}


func (s *TransactionService) GetAllTransactions(userID uint, page, limit int) ([]models.Transaction, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.Repo.GetAllTransactionsByUser(userID, limit, offset)
}
