package util

import (
	"e-meeting-api/configs"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = configs.AppConfig.JWTSecretKey

func GenerateToken(expiry time.Duration, claims jwt.MapClaims) (string, error) {
	claims["exp"] = time.Now().Add(expiry).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secret), nil
	})
}
