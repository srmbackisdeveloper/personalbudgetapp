package models

import (
	"time"
)

type AddCardRequest struct {
	CardNumber     string        `json:"cardNumber"`
	CardBalance    float64       `json:"cardBalance"`
	CardType       string        `json:"cardType"`
	CardExpireDate string     `json:"cardExpireDate"`
	AccountID      uint          `json:"accountId"`
}

func NewCard(number string, balance float64, _type string, expireDate time.Time, accountId uint) *Card {
	newCard := &Card{
		CardNumber:       number,
		CardBalance:    balance,
		CardType:   _type,
		CardExpireDate:    expireDate,
		AccountID:    accountId,
	}

	return newCard
}