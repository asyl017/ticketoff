package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

var jwtKey = []byte("your-secret-key") // Change this to your own secret key

func LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return

		}
		ctx := context.WithValue(r.Context(), "userID", claims.Issuer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
