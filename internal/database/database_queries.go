package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"personal_budget_app/internal/functionalities"
	"personal_budget_app/internal/models"
)

func (s *service) Health() map[string]string {
	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetAllAccounts() ([]*models.Account, error) {
	var accounts []*models.Account

	result := s.db.Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (s *service) GetAccount(id uint) (*models.Account, error) {
	var account models.Account

	result := s.db.Preload("Cards.OutgoingTransactions").Preload("Cards.IncomingTransactions").First(&account, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Use errors.Is for error comparison to properly handle wrapped errors.
			return nil, fmt.Errorf("account with id=%v not found", id)
		}
		return nil, result.Error
	}

	return &account, nil  // Return a pointer to the loaded account
}

func (s *service) CreateAccount(account *models.Account) error {
	result := s.db.Create(account)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Successfully created user with id: %v\n", account.ID)

	return nil
}

func (s *service) DeleteAccount(id uint) error {
	result := s.db.Delete(&models.Account{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id=%v is not found", id)
	}

	fmt.Printf("Successfully deleted user with id: %v\n", id)

	return nil
}

func (s *service) UpdateAccount(id uint, accountUpdates *models.UpdateAccountRequest) error {
	result := s.db.Model(&models.Account{}).Where("id = ?", id).Updates(accountUpdates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id=%v not found", id)
	}

	fmt.Printf("Successfully updated user with id: %v\n", id)

	return nil
}

func (s *service) AuthenticateUser(email, password string) (bool, error) {
	var account models.Account

	if err := s.db.Where("email = ?", email).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	result := functionalities.CheckPassword(account.Password, password)
	return result, nil
}



func (s *service) GetIdByEmail(email string) (uint, error)  {
	var account models.Account

	if err := s.db.Where("email = ?", email).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
		return 0, err
	}

	return account.ID, nil
}


