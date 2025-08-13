package repository

import (
	"gorm.io/gorm"
	"vaqua/models"
"time"	)

type TransactionRepository interface {
	FindByID(id uint, userID uint) (*models.Transaction, error)
	FindExpensesByUserAndPeriod(userID uint, period string) ([]models.Transaction, error)
}

type TransactionRepo struct{
	DB *gorm.DB
}

func (r *TransactionRepo) FindByID(id uint, userID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
} 

func (r *TransactionRepo) FindExpensesByUserAndPeriod(userID uint, period string) ([]models.Transaction, error) {
	var expenses []models.Transaction

	//start date based on period
	var startDate time.Time
	now := time.Now()

	switch period {
	case "last_month":
		startDate = now.AddDate(0, -1, 0)
	case "last_6_months":
		startDate = now.AddDate(0, -6, 0)
	case "last_year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0) // default to last month
	}

	err := r.DB.
		Where("user_id = ? AND type = ? AND created_at >= ?", userID, "expense", startDate).
		Find(&expenses).Error
	
	return expenses, err	
}    