package utils

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
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
	return token.SignedString(jwtKey)
}

func GenerateToken(email string) string {
	claims := &jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func ParseToken(tokenString string) (string, error) {
	logrus.Info("Parsing token for ", tokenString[0:10])
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Error("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		logrus.Error("Error parsing token: ", err)
		return "", err
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		logrus.Error("Error parsing token: ", err)
		return "", err
	} else if err != nil {
		logrus.Error("Error parsing token: ", err)
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}
	logrus.Error("Error parsing token: ", jwt.ErrTokenInvalidClaims)
	return "", jwt.ErrTokenInvalidClaims
}
