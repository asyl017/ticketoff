package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"ticketoff/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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
		return []byte("your-secret-key"), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Issuer, nil
}

var jwtKey = []byte("your-secret-key")

func VerifyEmailToken(tokenString string) (string, error) {

	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return jwtKey, nil

	})

	if err != nil {

		return "", err

	}

	if !token.Valid {

		return "", errors.New("invalid token")

	}

	if claims.ExpiresAt < time.Now().Unix() {

		return "", errors.New("token expired")

	}

	return claims.Subject, nil

}
