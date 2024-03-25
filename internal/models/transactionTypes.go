package models

import "time"

type AddTransactionRequest struct {
	TransactionAmount float64   `json:"transactionAmount"`
	FromCardID            uint      `json:"fromCardID"`
	ToCardNumber string      `json:"toCardNumber"`
}

func NewTransaction(amount float64, fromCardId uint, toCardID uint) *Transaction {
	newCard := &Transaction{
		TransactionTime: time.Now(),
		TransactionAmount: amount,
		FromCardID: fromCardId,
		ToCardID: toCardID,
	}

	return newCard
}
