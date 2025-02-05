package handler

import (
	"log"
	"net/http"
	"ticketoff/repositories"
	"ticketoff/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyUser(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	log.Println("VerifyUser function called")
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Verification token is missing", http.StatusBadRequest)
		return
	}

	email, err := utils.ParseToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	log.Printf("Email to verify: %s", email)

	userRepo := repositories.NewUserRepository(db)
	err = userRepo.VerifyUser(email)
	if err != nil {
		log.Printf("Failed to verify user: %v", err)
		http.Error(w, "Failed to verify user", http.StatusInternalServerError)
		return
	}

	log.Println("Email verified successfully")
	w.Write([]byte("Email verified successfully"))
}
