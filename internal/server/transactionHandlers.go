package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"personal_budget_app/internal/functionalities"
	"personal_budget_app/internal/models"
	"strconv"
)

func (s *Server) handleAddTransactionTo(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	// auth check passed:
	tsLimits := map[string]float64{
		"MIN_AMOUNT": 100.0,
		"MAX_AMOUNT": 100000.0,
	}

	req := new(models.AddTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "E: " + err.Error()})
		return
	}

	// get receiver ID
	toCardID, err := s.db.FindCardIDByCardNumber(req.ToCardNumber)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusNotFound, APIServerError{Error: "Destination card not found"})
		return
	}

	// check card belongs
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}



	doesBelong, err := s.db.CheckCardBelongsToUser(req.FromCardID, uint(userID))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	if !doesBelong {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: fmt.Sprintf("The card (id=%v) is private and does not belong to this user", req.FromCardID)})
		return
	}

	// check balance
	if req.TransactionAmount < tsLimits["MIN_AMOUNT"] { // MIN
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: fmt.Sprintf("Minimum transaction amount is %v", tsLimits["MIN_AMOUNT"])})
		return
	}

	if req.TransactionAmount > tsLimits["MAX_AMOUNT"] { // MAX
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: fmt.Sprintf("Maximum transaction amount is %v", tsLimits["MAX_AMOUNT"])})
		return
	}

	hasBalance, err := s.db.CheckCardBalance(req.FromCardID, req.TransactionAmount)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}
	if !hasBalance {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Insufficient balance"})
		return
	}

	// create the transaction
	ts := models.NewTransaction(req.TransactionAmount, req.FromCardID, toCardID)
	if err := s.db.AddTransaction(ts); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	// sender's account (-)
	if err := s.db.SenderUpdateBalance(req.FromCardID, req.TransactionAmount); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	// receiver's account (+)
	if err := s.db.ReceiverUpdateBalance(toCardID, req.TransactionAmount); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Transaction successful"})
}


func (s *Server) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	// auth check passed:
	idString := mux.Vars(r)["cardId"]
	cardId, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "invalid card id"})
		return
	}

	// check card belongs
	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	doesBelong, err := s.db.CheckCardBelongsToUser(uint(cardId), uint(userID))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	if !doesBelong {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: fmt.Sprintf("The card (id=%v) is private and does not belong to this user", cardId)})
		return
	}

	// main
	transactionType := r.URL.Query().Get("type")

	var transactions []*models.Transaction

	switch transactionType {
	case "incoming":
		transactions, err = s.db.GetIncomingTransactions(uint(cardId))
	case "outgoing":
		transactions, err = s.db.GetOutgoingTransactions(uint(cardId))
	default:
		transactions, err = s.db.GetAllTransactions(uint(cardId))
	}

	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, transactions)
}