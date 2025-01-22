package handler

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"net/http"
	"ticketoff/repositories"
	"ticketoff/utils"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type EmailHandler struct {
	UserRepo repositories.UserRepository
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	var emailReq EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		fmt.Println("Invalid request payload", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", Dialer.Username)
	m.SetHeader("To", emailReq.To)
	m.SetHeader("Subject", emailReq.Subject)
	m.SetBody("text/plain", emailReq.Message)

	if err := Dialer.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email", err)
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}

// ConfirmEmail handler
func (h *EmailHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, err := utils.ParseToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}
	err = h.UserRepo.ConfirmEmail(email)
	if err != nil {
		http.Error(w, "Error confirming email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email confirmed successfully!"))
}
