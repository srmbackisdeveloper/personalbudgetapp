package database

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"personal_budget_app/internal/models"
)

type Service interface {
	Health() map[string]string
	GetAllAccounts() ([]*models.Account, error)
	GetAccount(id uint) (*models.Account, error)
	CreateAccount(account *models.Account) error
	DeleteAccount(id uint) error
	UpdateAccount(id uint, accountUpdates *models.UpdateAccountRequest) error

	// AuthenticateUser auth
	AuthenticateUser(email, password string) (bool, error)
	GetIdByEmail(email string) (uint, error)

	// AddCard cards
	AddCard(card *models.Card) error
	DeleteCard(id uint) error
	GetCards(accountID uint) ([]*models.Card, error)
	GetCard(cardID uint) (*models.Card, error)

	// AddTransaction transaction
	AddTransaction(ts *models.Transaction) error
	CheckCardBalance(cardID uint, amount float64) (bool, error)
	CheckCardBelongsToUser(cardId, accountId uint) (bool, error)

	SenderUpdateBalance(cardID uint, amount float64) error
	ReceiverUpdateBalance(cardID uint, amount float64) error

	FindCardIDByCardNumber(cardNumber string) (uint, error)

	// GetAllTransactions get
	GetAllTransactions(cardId uint) ([]*models.Transaction, error)
	GetIncomingTransactions(cardId uint) ([]*models.Transaction, error)
	GetOutgoingTransactions(cardId uint) ([]*models.Transaction, error)

	// Settings
	SetDefaultCard(userId, cardId uint) (error)

	CreatePasswordResetToken(token models.PasswordResetToken) error
	ValidateToken(token string) (uint, error)
	UpdatePassword(accountID uint, newPassword string) error
	MarkTokenAsUsed(token string) error

	CheckCurrentPassword(accountID uint, currentPassword string) (bool, error)
}

type service struct {
	db *gorm.DB
}

func New() Service {
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port     := os.Getenv("DB_PORT")
	host     := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate models
	err = db.AutoMigrate(&models.Account{}, &models.Card{}, &models.Transaction{}, &models.PasswordResetToken{})
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	return &service{db: db}
}

