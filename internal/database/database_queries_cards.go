package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"personal_budget_app/internal/models"
)

func (s *service) AddCard(card *models.Card) error {
	// Check the current number of cards for the account
	var count int64
	result := s.db.Model(&models.Card{}).Where("account_id = ?", card.AccountID).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	// Check if the account has already reached the maximum number of cards
	if count >= 3 {
		return fmt.Errorf("maximum number of cards (3) for account (id=%v) reached", card.AccountID)
	}

	// If not, proceed with adding the new card
	result = s.db.Create(card)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Successfully created card (id=%v) for user (id=%v)\n", card.ID, card.AccountID)
	return nil
}

func (s *service) DeleteCard(id uint) error {
	result := s.db.Delete(&models.Card{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("card with id=%v is not found", id)
	}

	fmt.Printf("Successfully deleted card with id: %v\n", id)

	return nil
}


func (s *service) GetCards(accountID uint) ([]*models.Card, error) {
	var cards []*models.Card

	err := s.db.Where("account_id = ?", accountID).
		Preload("OutgoingTransactions").
		Preload("IncomingTransactions").
		Find(&cards).Error

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *service) GetCard(cardID uint) (*models.Card, error) {
	var card *models.Card  // Use a non-pointer Card struct here to avoid nil pointer dereference issues

	err := s.db.Preload("OutgoingTransactions").Preload("IncomingTransactions").First(&card, cardID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("card with ID %d not found", cardID)
		}
		return nil, err
	}

	return card, nil
}
