package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"vaqua/models"
	"vaqua/repository"
)

type TransferService interface {
	TransferFunds(userID uint, req *models.TransferRequest) error
}

type TransferServices struct {
	Repo repository.TransferRepo
}

//connection to transfer repo

/*
function should:
get logged in user account with is authenticated from the context
get the recipient's account by account number
to not transfer to self
check if user has sufficient balance
update account in the database for both sender and recipient
create a transaction record for the transfer
*/

func (s *TransferServices) TransferFunds(senderUserID uint, request *models.TransferRequest) error {

	// Get logged in user account which is the sender account

	senderAcc, err := s.Repo.GetAccountByUserID(senderUserID)
	if err != nil {
		return errors.New("sender account not found")
	}

	// recipient account by account number

	recipientAcc, err := s.Repo.GetAccountByAccountNumber(request.RecipientAcc)
	if err != nil {
		return errors.New("recipient account not found")
	}

	// not allow transfer to self

	if senderAcc.AccountNumber == recipientAcc.AccountNumber {
		return errors.New("cannot transfer to your own account")
	}

	// check users balance, if user balance is less than the amount they want to transfer error here
	if senderAcc.Balance < request.Amount {
		return errors.New("insufficient balance")
	}

	fmt.Printf("Transferring %.2f from %s to %s\n", request.Amount, senderAcc.AccountNumber, recipientAcc.AccountNumber)

	// Save updated recipient account balance via atomic transaction logic
	if err := s.executeTransferTransaction(senderAcc, recipientAcc, request); err != nil {
		return err
	}

	return nil
}

//atomic transfer transaction
// actual transfer logic is handled here via a transaction repo method
// here both the sender and recipient account updates are atomic compared to normal service one which would update the sender account and then the recipient account
//here if one fails, it rolls back the changes made to both accounts

func (s *TransferServices) executeTransferTransaction(senderAcc, recipientAcc *models.Account, request *models.TransferRequest) error {
	return s.Repo.WithTransaction(func(txRepo *repository.TransferRepo, tx *gorm.DB) error { //with transaction method in repo layer

		//update balances
		senderAcc.Balance -= request.Amount
		recipientAcc.Balance += request.Amount

		if err := txRepo.UpdateAccount(senderAcc, tx); err != nil {
			return errors.New("failed to update sender account balance")
		}
		if err := txRepo.UpdateAccount(recipientAcc, tx); err != nil {
			return errors.New("failed to update recipient account balance")
		}

		if err := txRepo.CreateTransfer(senderAcc.UserID, recipientAcc.UserID, request.Amount, request.Description, "expense", tx); err != nil {
			return fmt.Errorf("failed to create sender transaction: %w", err)
		}
		
		if err := txRepo.CreateTransfer(recipientAcc.UserID, senderAcc.UserID, request.Amount, request.Description, "income", tx); err != nil {
			return fmt.Errorf("failed to create recipient transaction: %w", err)
		}
		return nil
	})
}

//helper to create transaction
