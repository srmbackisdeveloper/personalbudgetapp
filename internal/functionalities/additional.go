package functionalities

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
)

func WriteJSON(w http.ResponseWriter, status int, anything interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(anything)
	if err != nil {
		return
	}
}

func GenerateSecureToken() string {
	const tokenSize = 32

	tokenBytes := make([]byte, tokenSize)

	if _, err := rand.Read(tokenBytes); err != nil {
		log.Printf("Failed to generate secure token: %v", err)
		return ""
	}

	return hex.EncodeToString(tokenBytes)
}

func SendEmail(email, link string) error  {
	fromEmail := os.Getenv("FROM_EMAIL")
	fromEmailSecret := os.Getenv("FROM_EMAIL_PASSWORD")

	m := gomail.NewMessage()

	m.SetHeader("From", fromEmail)
	m.SetHeader("To", email)

	m.SetHeader("Subject", "Password Recovery - Personal Budget App")
	m.SetBody("text/html", fmt.Sprintf("Here is your link to recover the password:<br>%v<br><br>Warning: token expires after 10 minutes!", link))

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, fromEmailSecret)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}