package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"ticketoff/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your-secret-key")

func GenerateJWT(user *models.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    fmt.Sprintf("%d", user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func GenerateToken(email string) string {
	token := make([]byte, 16)
	rand.Read(token)
	return hex.EncodeToString(token)
}

func ParseToken(tokenString string) (string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Subject, nil
}
