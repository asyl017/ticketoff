package handler

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ticketoff/models"
	"ticketoff/repositories"
	"ticketoff/utils"
	"time"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	userRepo repositories.UserRepository
}

func NewAuthHandler(userRepo repositories.UserRepository) AuthHandler {
	return &authHandler{userRepo: userRepo}
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.Logger.WithError(err).Error("Invalid request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	utils.Logger.Info("Logging in user with email: ", creds.Email)
	user, err := a.userRepo.GetUserByEmail(creds.Email)
	if err != nil {
		utils.Logger.WithError(err).Error("Invalid email or password")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	utils.Logger.Info("User found: ", user, ". Comparing passwords")
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		utils.Logger.WithError(err).Error("Invalid email or password")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	utils.Logger.Info("Password is correct. Generating token")
	token, err := generateJWT(user)
	if err != nil {
		utils.Logger.WithError(err).Error("Error generating token")
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	utils.Logger.Info("Token generated successfully. Sending response")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func generateJWT(user *models.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    string(user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}
