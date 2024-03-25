package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"personal_budget_app/internal/models"
)

func (s *service) AddTransaction(ts *models.Transaction) error {
	var card models.Card
	result := s.db.Where("id = ?", ts.ToCardID).First(&card)

	if result.Error != nil {
		fmt.Printf("Error: Destination card (id=%v) does not exist: %v\n", ts.ToCardID, result.Error)
		return result.Error
	}


	result = s.db.Create(ts)
	if result.Error != nil {
		fmt.Printf("Error adding transaction for card (id=%v): %v\n", ts.FromCardID, result.Error)
		return result.Error
	}

	fmt.Printf("Successfully added transaction (id=%v): [%v --> %v];\n", ts.ID, ts.FromCardID, ts.ToCardID)

	return nil
}

func (s *service) FindCardIDByCardNumber(cardNumber string) (uint, error) {
	var card models.Card
	result := s.db.Where("card_number = ?", cardNumber).First(&card)
	if result.Error != nil {
		return 0, result.Error
	}
	return card.ID, nil
}



func (s *service) CheckCardBalance(cardID uint, amount float64) (bool, error) {
	var card models.Card

	// Retrieve the card by ID
	result := s.db.First(&card, "id = ?", cardID)
	if result.Error != nil {
		fmt.Printf("Error retrieving card (id=%v): %v\n", cardID, result.Error)
		return false, result.Error
	}

	// Check if the card has enough balance
	if card.CardBalance < amount {
		return false, nil // Not enough balance, but no error occurred
	}

	return true, nil // Sufficient balance
}

func (s *service) SenderUpdateBalance(cardID uint, amount float64) error {
	var card models.Card

	// Retrieve the card by ID
	result := s.db.First(&card, "id = ?", cardID)
	if result.Error != nil {
		fmt.Printf("Error retrieving card (id=%v): %v\n", cardID, result.Error)
		return result.Error
	}

	card.CardBalance -= amount

	updateResult := s.db.Save(&card)
	if updateResult.Error != nil {
		fmt.Printf("Error updating card balance (id=%v): %v\n", cardID, updateResult.Error)
		return updateResult.Error
	}

	fmt.Printf("Successfully updated card balance (id=%v)\n", cardID)
	return nil
}

func (s *service) ReceiverUpdateBalance(cardID uint, amount float64) error {
	var card models.Card

	// Retrieve the card by ID
	result := s.db.First(&card, "id = ?", cardID)
	if result.Error != nil {
		fmt.Printf("Error retrieving card (id=%v): %v\n", cardID, result.Error)
		return result.Error
	}

	card.CardBalance += amount

	updateResult := s.db.Save(&card)
	if updateResult.Error != nil {
		fmt.Printf("Error updating card balance (id=%v): %v\n", cardID, updateResult.Error)
		return updateResult.Error
	}

	fmt.Printf("Successfully updated card balance (id=%v)\n", cardID)
	return nil
}

func (s *service) CheckCardBelongsToUser(cardId, accountId uint) (bool, error) {
	var card models.Card
	result := s.db.First(&card, cardId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	if card.AccountID == accountId {
		return true, nil // Card belongs to the user
	}

	return false, nil
}



// get
func (s *service) GetIncomingTransactions(cardId uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	// Query for transactions where the card is the recipient
	result := s.db.Where("to_card_id = ?", cardId).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (s *service) GetOutgoingTransactions(cardId uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	// Query for transactions where the card is the sender
	result := s.db.Where("from_card_id = ?", cardId).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (s *service) GetAllTransactions(cardId uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	// Query for all transactions related to the card, either as sender or recipient
	result := s.db.Where("from_card_id = ? OR to_card_id = ?", cardId, cardId).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}
