package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", s.handleLogin).Methods("POST")
	router.HandleFunc("/register", s.handleCreateAccount).Methods("POST")

	// Recover pass
	router.HandleFunc("/accounts/forgetpw", s.handleForgetPassword).Methods("POST")

	router.HandleFunc("/accounts/reset-password", s.handlePasswordReset).Methods("POST")


	// Protected routes
	secure := router.PathPrefix("/api").Subrouter()
	// secure.Use(requestLoggerMiddleware) // for debug
	secure.Use(JWTMiddleware)


	secure.HandleFunc("/logout", s.handleLogout).Methods("GET", "POST")
	secure.HandleFunc("/accounts", s.handleGetAccounts).Methods("GET")
	secure.HandleFunc("/accounts/{id}", s.handleGetAccount).Methods("GET")
	secure.HandleFunc("/accounts/{id}", s.handleDeleteAccount).Methods("DELETE")
	secure.HandleFunc("/accounts/{id}", s.handleUpdateAccount).Methods("PUT")

	secure.HandleFunc("/cards", s.handleAddCard).Methods("POST")
	secure.HandleFunc("/cards", s.handleGetCards).Methods("GET")
	secure.HandleFunc("/cards/{id}", s.handleDeleteCard).Methods("DELETE")
	secure.HandleFunc("/cards/{id}", s.handleGetCard).Methods("GET")

	secure.HandleFunc("/transaction/{cardId}", s.handleGetTransactions).Methods("GET")
	secure.HandleFunc("/transaction", s.handleAddTransactionTo).Methods("POST")


	// account settings
	secure.HandleFunc("/accounts/settings/default-card/{cardId}", s.handleSetDefaultCard).Methods("POST")
	secure.HandleFunc("/accounts/settings/change-password/{id}", s.handleUpdatePassword).Methods("PUT")

	corsRouter := corsMiddleware(router)

	return corsRouter
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}