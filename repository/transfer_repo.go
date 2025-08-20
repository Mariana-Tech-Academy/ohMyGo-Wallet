package repository

import (
	"errors"
	"vaqua/db"
	"vaqua/models"

	"gorm.io/gorm"
)


//ammend code to reflect the recipient's account number

type TransferRepository interface {
	GetAccountByUserID(userID uint) (*models.Account, error)
	GetAccountByAccountNumber(accountNumber string) (*models.Account, error)
	UpdateAccount(account *models.Account) error
	CreateTransfer(senderAcc, recipientAcc uint, amount float64, description string, tx *gorm.DB) error
	CreateAccount(account *models.Account) error
	WithTransaction(fn func(txRepo *TransferRepo, tx *gorm.DB) error) error

}

type TransferRepo struct{
	DB *gorm.DB
}


func NewTransferRepo() *TransferRepo {
	return &TransferRepo{DB : db.Db }
}  // NewTransferRepo creates a new instance of TransferRepo which im using for when i do the create account for the user during signup

func (r *TransferRepo) GetAccountByUserID(userID uint) (*models.Account, error) {
	var account models.Account

	err := db.Db.Where("user_id = ?", userID).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No account found for this user - handle accordingly
		return nil, errors.New("user account not found")
	} else if err != nil {
		// Other errors
		return nil, err
	}
	return &account, nil
}


// GetAccountByAccountNumber retrieves an account by its account number so we can find the recipient's account during a transfer
func (r *TransferRepo) GetAccountByAccountNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	err := db.Db.Where("account_number = ?", accountNumber).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("recipient account not found")
	} else if err != nil {
		return nil, err
	}
	return &account, nil
}

//we need to pass a transaction object if it is used in a transaction


// we need to update the account after the transfer so it can reflect the new balance, it gets saved to the database
func (r *TransferRepo) UpdateAccount(account *models.Account, tx *gorm.DB) error {
	if tx != nil {
		return tx.Save(account).Error
	}
	return db.Db.Save(account).Error
}


//We create transfer in the database
func (r *TransferRepo) CreateTransfer(senderAcc, recipientAcc uint, amount float64, description, txType string, tx *gorm.DB) error {
	transfer := models.Transaction{//create a new transfer record
		UserID:      senderAcc,
		RecipientID: recipientAcc,
		Type:        txType, //updated
		Amount:      amount,
		Description: description,
	}
	if tx != nil {
		return tx.Create(&transfer).Error
	}
	return db.Db.Create(&transfer).Error
}

//is needed to create an account for the user during signup linked to the userrepo accounts so left it back here for that
func (r *TransferRepo) CreateAccount(account *models.Account) error {
	return db.Db.Create(account).Error
} 


//withtransaction helper is used to wrap the create and update operations in a transaction
//the context is passed to the function so it can be used within the transaction
//function takes other function as parameter this way, it can be reused for different operations with the database
func (r *TransferRepo) WithTransaction(fn func(txRepo *TransferRepo, tx *gorm.DB) error) error {
	return db.Db.Transaction(func(tx *gorm.DB) error {
		// Create a new repo bound to THIS transaction
		txRepo := &TransferRepo{DB : tx}
		return fn(txRepo, tx)
	})
}

