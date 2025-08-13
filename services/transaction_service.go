package services

import ("vaqua/repository"
		"vaqua/models"
		"errors")


type TransactionService struct {
	Repo repository.TransactionRepository
}

func (s *TransactionService) GetTransactionByID(id uint, userID uint) (*models.Transaction, error) {
	transaction, err := s.Repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	//redunduncy check needed?
	if transaction.UserID != userID {
		return nil, errors.New("unauthorized access")

	}
	return transaction, nil
}

func (s *TransactionService) GetExpensesByUserAndPeriod(userID uint, period string) ([]models.Transaction, error) {
	// TODO: Implement logic to retrieve expenses based on the period and UserID
	return s.Repo.FindExpensesByUserAndPeriod(userID, period)
}

