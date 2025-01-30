package middleware

import (
	"net/http"
	"ticketoff/repositories"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your-secret-key") // Change this to your own secret key

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type AuthHandler struct {
	UserRepo repositories.UserRepository
}

func NewAuthHandler(userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

func (h *AuthHandler) Authorize(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			user, err := h.UserRepo.GetUserByEmail(claims.Email)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
