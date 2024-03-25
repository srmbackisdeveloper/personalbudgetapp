package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"personal_budget_app/internal/functionalities"
	"personal_budget_app/internal/models"
	"strconv"
	"time"
)

func (s *Server) handleAddCard(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	// after auth

	id, err := strconv.Atoi(user.UserID)

	req := new(models.AddCardRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	var expireDate time.Time
	if req.CardExpireDate != "" {
		var parseErr error
		expireDate, parseErr = time.Parse("2006-01-02", req.CardExpireDate) // Added req.CardExpireDate as argument
		if parseErr != nil {
			functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid date format: " + parseErr.Error()})
			return
		}
	}

	card := models.NewCard(req.CardNumber, req.CardBalance, req.CardType, expireDate, uint(id))

	err = s.db.AddCard(card)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, card)
}


func (s *Server) handleDeleteCard(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]

	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	// after the verification

	idCard, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "invalid id"})
		return
	}

	// check card belongs
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	doesBelong, err := s.db.CheckCardBelongsToUser(uint(idCard), uint(userID))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	if !doesBelong {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: fmt.Sprintf("The card (id=%v) is private and does not belong to this user", idCard)})
		return
	}

	//

	err = s.db.DeleteCard(uint(idCard))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}


	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "card successfully deleted"})
}

func (s *Server) handleGetCards(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	// after the verification

	accountId, err := strconv.Atoi(user.UserID)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	cards, err := s.db.GetCards(uint(accountId))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}


	functionalities.WriteJSON(w, http.StatusOK, cards)
}


func (s *Server) handleGetCard(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	// after the verification

	idCard, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "invalid id"})
		return
	}

	// check card belongs
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	//if idCard == 0 {
	//	cardZero := new(models.Card)
	//	functionalities.WriteJSON(w, http.StatusOK, cardZero)
	//	return
	//}

	//

	doesBelong, err := s.db.CheckCardBelongsToUser(uint(idCard), uint(userID))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	if !doesBelong {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: fmt.Sprintf("The card (id=%v) is private and does not belong to this user", idCard)})
		return
	}

	//

	card, err := s.db.GetCard(uint(idCard))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "could not fetch a card"})
		return
	}


	functionalities.WriteJSON(w, http.StatusOK, card)
}
