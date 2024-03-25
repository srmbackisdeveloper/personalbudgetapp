package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"personal_budget_app/internal/database"
)

type APIServerError struct {
	Error string `json:"error"`
}

type Server struct {
	port int
	db database.Service
	sessionStore *sessions.CookieStore
}

func NewServer() *http.Server {
	err := godotenv.Load("../.././.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	sessionKey := []byte("secret") // !!! CONTINUE WORKING WITH sessions

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db: database.New(),
		sessionStore: sessions.NewCookieStore(sessionKey),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("server running on port: %v\n", port)
	return server
}