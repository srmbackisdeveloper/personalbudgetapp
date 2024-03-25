package server

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"personal_budget_app/internal/functionalities"
	"strconv"
	"time"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}


// Login
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromToken(r)
	if err == nil && user.UserID != "" {
		functionalities.WriteJSON(w, http.StatusForbidden, APIServerError{Error: "Already logged in"})
		return
	}


	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid request payload"})
		return
	}

	log.Printf("Login attempt for email: %s", loginRequest.Email)

	authenticated, err := s.db.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "Error during authentication"})
		return
	}

	if !authenticated {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Invalid credentials"})
		return
	}

	userID, err := s.db.GetIdByEmail(loginRequest.Email)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}

	claims := &Claims{
		UserID: strconv.Itoa(int(userID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("syrymbek"))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "Failed to generate token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().UTC().Add(24 * time.Hour),
		HttpOnly: false, // Prevents client-side JS from accessing the cookie
		Secure:   false, // Ensures cookie is sent over HTTPS
		Path:     "/",  // Cookie available to entire domain
	})


	log.Printf("SUCCESS: %s;", loginRequest.Email)
	functionalities.WriteJSON(w, http.StatusOK, tokenString)
}


// Logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	})

	log.Printf("Logout SUCCESS;")
	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}





func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "dashboard get success"})
}
