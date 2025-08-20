package repository

import (
	"time"
	"vaqua/db"
	"vaqua/models"


)

type TransactionRepo struct{}

type TransactionRepository interface {
	GetIncomeByPeriod(userID uint, start, end time.Time) ([]models.Transaction, error)
	GetExpensesByPeriod(userID uint, start, end time.Time) ([]models.Transaction, error)
	GetUserBalanceByID(userID uint) (float64, error)
	GetAllTransactionsByUser(userID uint, limit int, offset int) ([]models.Transaction, error)
	GetTransactionsByUserID(userID uint) ([]models.Transaction, error)
}



func (r *TransactionRepo) GetIncomeByPeriod(userID uint, start, end time.Time) ([]models.Transaction, error) {
    var incomes []models.Transaction
    err := db.Db.Where("user_id = ? AND type = ? AND created_at BETWEEN ? AND ?", userID, "income", start, end).
        Find(&incomes).Error

    return incomes, err
}

func (r *TransactionRepo) GetExpensesByPeriod(userID uint, start, end time.Time) ([]models.Transaction, error) {
	var expenses []models.Transaction
	err := db.Db.Where("user_id = ? AND type = ? AND created_at BETWEEN ? AND ?", userID, "expenses", start, end).
	Find(&expenses).Error

	return expenses, err

}


func (r *TransactionRepo) GetUserBalanceByID(userID uint) (float64, error) {
    var user models.Account
    err := db.Db.Select("balance").Where("id = ?", userID).First(&user).Error
    if err != nil {
        return 0, err
    }
    return user.Balance, nil
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

func (r *TransactionRepo) GetTransactionByUserID(userID uint) (models.Transaction, error) {
	var transaction models.Transaction
	err := db.Db.Where("user_id", userID).Find(&transaction).Error
	if err!= nil {
		return models.Transaction{}, err
	}
	return transaction, nil 
}


