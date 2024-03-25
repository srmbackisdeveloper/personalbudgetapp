package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"personal_budget_app/internal/functionalities"
	"personal_budget_app/internal/models"
	"strconv"
	"time"
)


func ExtractUserFromToken(r *http.Request) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_TOKEN")

	c, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	tokenString := c.Value

	if tokenString == "" {
		return nil, fmt.Errorf("token not found in cookies")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}



// GET ALL USERS
func (s *Server) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.db.GetAllAccounts()
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, accounts)
}

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]

	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}

	if user.UserID != idString {
		functionalities.WriteJSON(w, http.StatusForbidden, APIServerError{Error: "Access denied"})
		return
	}

	// after the verification

	id, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "invalid id"})
		return
	}

	account, err := s.db.GetAccount(uint(id))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, account)
}

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	createAccReq := new(models.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	hashedPassword, err := functionalities.HashPassword(createAccReq.Password)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	var birthday time.Time
	if createAccReq.Birthday != "" {
		var parseErr error
		birthday, parseErr = time.Parse("2006-01-02", createAccReq.Birthday)
		if parseErr != nil {
			functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid birthday format: " + parseErr.Error()})
			return
		}
	}

	account := models.NewAccount(
		createAccReq.Email,
		hashedPassword,
		createAccReq.FirstName,
		createAccReq.LastName,
		birthday,
		createAccReq.PhoneNumber,
	)

	err = s.db.CreateAccount(account)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, account)
}


func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]

	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}

	if user.UserID != idString {
		functionalities.WriteJSON(w, http.StatusForbidden, APIServerError{Error: "Access denied"})
		return
	}

	// after the verification

	id, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "invalid id"})
		return
	}

	err = s.db.DeleteAccount(uint(id))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}


	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "successfully deleted"})
}


func (s *Server) handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "invalid id"})
		return
	}

	updateAccReq := new(models.UpdateAccountRequest)
	if err = json.NewDecoder(r.Body).Decode(&updateAccReq); err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "invalid request body: " + err.Error()})
		return
	}

	err = s.db.UpdateAccount(uint(id), updateAccReq)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Account successfully updated"})
}


