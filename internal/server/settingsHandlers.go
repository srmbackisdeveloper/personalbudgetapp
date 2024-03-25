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



func (s *Server) handleSetDefaultCard(w http.ResponseWriter, r *http.Request) {
	user, err := ExtractUserFromToken(r)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: err.Error()})
		return
	}
	if user.UserID == "" {
		functionalities.WriteJSON(w, http.StatusUnauthorized, APIServerError{Error: "Unauthorized"})
		return
	}

	userIdString := user.UserID
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid card id"})
		return
	}

	// after auth:
	vars := mux.Vars(r)
	cardIdString := vars["cardId"]
	cardId, err := strconv.Atoi(cardIdString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid card id"})
		return
	}

	// check card belongs
	if cardId != 0 {
		doesBelong, err := s.db.CheckCardBelongsToUser(uint(cardId), uint(userId))
		if err != nil {
			functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
			return
		}

		if !doesBelong {
			functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: fmt.Sprintf("The card (id=%v) is private and does not belong to this user", cardId)})
			return
		}
	}

	err = s.db.SetDefaultCard(uint(userId), uint(cardId))
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Default card id set successfully"})
}

// FORGET
func (s *Server) handleForgetPassword(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid request"})
		return
	}

	token := functionalities.GenerateSecureToken()

	accountID, err := s.db.GetIdByEmail(requestData.Email)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Email not registered"})
		return
	}

	// Create a password reset token record
	err = s.db.CreatePasswordResetToken(models.PasswordResetToken{
		AccountID: accountID,
		Token:     token,
	})
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "Failed to process password reset request"})
		return
	}

	// Construct the secure password recovery link
	link := fmt.Sprintf("http://localhost:3000/recovery/%s", token)

	// Send the recovery link to the user's email
	if err = functionalities.SendEmail(requestData.Email, link); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "Error sending email"})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Recovery link has been sent to your email successfully"})
}



// RECOVER
func (s *Server) handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Token       string `json:"token"`
		NewPassword string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid request format"})
		return
	}

	accountID, err := s.db.ValidateToken(requestData.Token)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Invalid or expired token"})
		return
	}

	if err := s.db.UpdatePassword(accountID, requestData.NewPassword); err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: "Failed to update password"})
		return
	}

	if err := s.db.MarkTokenAsUsed(requestData.Token); err != nil {
		// Log the error but don't fail the request; the password has been updated successfully
	}

	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Your password has been reset successfully."})
}

// update password
func (s *Server) handleUpdatePassword(w http.ResponseWriter, r *http.Request) {
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

	// after

	id, err := strconv.Atoi(idString)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "invalid id"})
		return
	}

	var updatePassReq struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err = json.NewDecoder(r.Body).Decode(&updatePassReq); err != nil {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "invalid request body: " + err.Error()})
		return
	}

	//log.Printf("Current: %s\n", updatePassReq.CurrentPassword)
	//log.Printf("New: %s\n", updatePassReq.NewPassword)

	if updatePassReq.CurrentPassword == updatePassReq.NewPassword {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Passwords should differ, try again"})
		return
	}

	passwordMatches, err := s.db.CheckCurrentPassword(uint(id), updatePassReq.CurrentPassword)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	if !passwordMatches {
		functionalities.WriteJSON(w, http.StatusBadRequest, APIServerError{Error: "Current password does not match, try again"})
		return
	}


	err = s.db.UpdatePassword(uint(id), updatePassReq.NewPassword)
	if err != nil {
		functionalities.WriteJSON(w, http.StatusInternalServerError, APIServerError{Error: err.Error()})
		return
	}

	functionalities.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password is successfully updated"})
}





