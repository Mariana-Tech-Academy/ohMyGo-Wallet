package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Firstname      *string `json:"firstname,omitempty"`
	Lastname       *string `json:"lastname,omitempty"`
	Email          string  `json:"email" gorm:"unique;not null"`
	Password       string  `json:"password" gorm:"not null"`
	AccountNumber   string `json:"account_number" gorm:"unique;not null"`
	Phonenumber    *uint `json:"phone_number,omitempty"`
	//uint consistency throughout the app
}

//moved account balance to account stuct

type Account struct {
	gorm.Model
	UserID        uint     `json:"-" gorm:"not null"` // Foreign key to User
	AccountNumber string  `json:"account_number" gorm:"uniqueIndex;not null"` // is a string
	Balance       float64  `json:"balance" gorm:"not null;default:0"`
	Phonenumber    *string `json:"phone_number,omitempty"`
	CurrentBalance float64 `json:"current_balance" gorm:"default:0"`

}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Firstname   *string `json:"firstname" binding:"required"`
	Lastname    *string `json:"lastname" binding:"required"`
	Phonenumber *uint `json:"phone_number" binding:"required,min=10000000000,max=99999999999999"`

}

