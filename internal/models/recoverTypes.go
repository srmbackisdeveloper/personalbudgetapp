package models

import "time"

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	AccountID    uint      `gorm:"not null"`
	Token     string    `gorm:"not null;size:255"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
}