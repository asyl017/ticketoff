package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"ticketoff/models"
	"ticketoff/repositories"
	"ticketoff/utils"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo repositories.UserRepository
}

func NewAuthHandler(userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(creds.Email)
	if err != nil {
		utils.Logger.Info("User not found")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		utils.Logger.WithFields(logrus.Fields{
			"email":    creds.Email,
			"password": creds.Password,
			"hash":     user.Password,
		}).Info("Password is incorrect")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	log.Println(user.EmailConfirmed)
	if !user.EmailConfirmed {
		http.Error(w, "Email not verified", http.StatusUnauthorized)
		return
	}

	token := utils.GenerateToken(user.Email)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Path:    "/",
		Value:   token,
		Expires: time.Now().Add(72 * time.Hour),
	})
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

/*func (h *AuthHandler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}*/
