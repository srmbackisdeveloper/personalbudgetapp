package models

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Email       string    `json:"email" gorm:"unique"`
	Password    string    `json:"-"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Birthday    time.Time `json:"birthday"`
	PhoneNumber string    `json:"phoneNumber"`
	DefaultCardID   uint      `json:"defaultCardID"` // I want to add here default card id
	Cards       []Card    `gorm:"foreignKey:AccountID" json:"cards,omitempty"`
}

type Card struct {
	gorm.Model
	CardNumber     string        `json:"cardNumber" gorm:"unique"`
	CardBalance    float64       `json:"cardBalance"`
	CardType       string        `json:"cardType"`
	CardExpireDate time.Time     `json:"cardExpireDate"`
	AccountID      uint          `json:"-"`
	OutgoingTransactions []Transaction `gorm:"foreignKey:FromCardID;references:ID" json:"outgoingTransactions,omitempty"`
	IncomingTransactions []Transaction `gorm:"foreignKey:ToCardID;references:ID" json:"incomingTransactions,omitempty"`
}

type Transaction struct {
	gorm.Model
	TransactionTime   time.Time `json:"transactionTime"`
	TransactionAmount float64   `json:"transactionAmount"`
	FromCardID            uint      `json:"fromCardID"`
	ToCardID   uint      `json:"toCardID"`
}

// ----------------------------------------

type CreateAccountRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Birthday    string `json:"birthday"`
	PhoneNumber string    `json:"phoneNumber"`
}

// ----------------------------------------
// Account things:

type UpdateAccountRequest struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
}

func NewAccount(email, password, firstName, lastName string, birthday time.Time, phoneNumber string) *Account {
	newAcc := &Account{
		Email:       email,
		Password:    password,
		FirstName:   firstName,
		LastName:    lastName,
		Birthday:    birthday,
		PhoneNumber: phoneNumber,
	}

	return newAcc
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

