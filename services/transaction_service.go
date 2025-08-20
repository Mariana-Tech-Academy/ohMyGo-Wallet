package services

import (
	"time"
	"vaqua/models"
	"vaqua/repository"
	

)
type TransactionServices interface {
	GetIncomeByPeriod(userID uint, start, end time.Time) ([]models.Transaction, float64, error) 
	GetExpensesByPeriod(userID uint, start, end time.Time) ([]models.Transaction, float64, error)
	GetAllTransactions(userID uint, page, limit int) ([]models.Transaction, error)
	GetTransactionByUserID(userID uint) (*models.Transaction, error)
    GetUserBalance(userID uint) (float64, error)

}

type TransactionService struct {
	Repo *repository.TransactionRepo
}


func NewTransactionService(repo *repository.TransactionRepo) *TransactionService {
	return &TransactionService{Repo: repo}
}


func (s *TransactionService) GetIncomeByPeriod(userID uint, start, end time.Time) ([]models.Transaction, float64, error) {
	income, err := s.Repo.GetIncomeByPeriod(userID, start, end)
	if err != nil {
		return nil, 0, err
	}

	var total float64
	for _, t := range income {
		total += t.Amount
	}
	return income, total, nil
}

func (s *TransactionService) GetExpensesByPeriod(userID uint, start, end time.Time) ([]models.Transaction, float64, error) {
	expenses, err := s.Repo.GetExpensesByPeriod(userID, start, end)
	if err != nil {
		return nil, 0, err
	}

	var total float64
	for _, t := range expenses {
		total += t.Amount
	}
	return expenses, total, nil
}

func (s *TransactionService) GetUserBalance(userID uint) (float64, error) {
    balance, err := s.Repo.GetUserBalanceByID(userID)
    if err != nil {
        return 0, err 
    }

    return balance, nil
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

func (s *TransactionService) GetTransactionByUserID(userID uint) (*models.Transaction, error) {
	transaction, err := s.Repo.GetTransactionByUserID(userID)
	if err != nil {
		return &models.Transaction{}, err
	}
	return &transaction, nil 

}

