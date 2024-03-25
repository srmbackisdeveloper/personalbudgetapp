package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"personal_budget_app/internal/functionalities"
	"personal_budget_app/internal/models"
	"time"
)

func (s *service) SetDefaultCard(userId, cardId uint) (error) {
	var account *models.Account
	result := s.db.Model(&account).Where("id = ?", userId).Update("default_card_id", cardId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *service) CreatePasswordResetToken(token models.PasswordResetToken) error {
	token.Used = false
	token.ExpiresAt = time.Now().Add(10 * time.Minute)

	result := s.db.Create(&token)
	if result.Error != nil {
		return result.Error
	}

	return nil
}


///
///
///

func (s *service) ValidateToken(token string) (uint, error) {
	var resetToken models.PasswordResetToken

	result := s.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&resetToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("invalid or expired token")
		}
		return 0, result.Error
	}

	return resetToken.AccountID, nil
}

func (s *service) UpdatePassword(accountID uint, newPassword string) error {
	hashedPassword, err := functionalities.HashPassword(newPassword)
	if err != nil {
		return err
	}

	result := s.db.Model(&models.Account{}).Where("id = ?", accountID).Update("password", string(hashedPassword))
	return result.Error
}


func (s *service) MarkTokenAsUsed(token string) error {
	result := s.db.Model(&models.PasswordResetToken{}).Where("token = ?", token).Update("used", true)
	return result.Error
}

func (s *service) CheckCurrentPassword(accountID uint, currentPassword string) (bool, error) {
	var account models.Account

	if err := s.db.First(&account, accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	result := functionalities.CheckPassword(account.Password, currentPassword)
	return result, nil
}
