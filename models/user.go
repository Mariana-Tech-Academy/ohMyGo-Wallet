package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Firstname      *string `json:"firstname,omitempty"`
	Lastname       *string `json:"lastname,omitempty"`
	Email          string  `json:"email" gorm:"unique;not null"`
	Password       string  `json:"password" gorm:"not null"`
	AccountNumber  uint64  `json:"account_number" gorm:"uniqueIndex;not null"`
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
	Phonenumber *string `json:"phone_number" binding:"required,len=11"`
}

