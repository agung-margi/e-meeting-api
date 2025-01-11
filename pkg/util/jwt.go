package util

import (
	"e-meeting-api/configs"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = configs.AppConfig.JWTSecretKey

func GenerateToken(secret string, expiry time.Duration, claims jwt.MapClaims) (string, error) {
	claims["exp"] = time.Now().Add(expiry).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv(secret)))
}

func VerifyToken(secret string, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(os.Getenv(secret)), nil
	})
}

func GeneratePasswordResetToken(userID int) (string, error) {
	// Membuat token dengan expiration time 30 menit
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Membuat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "",
			fmt.Errorf("failed to sign the token: %w", err)
	}

	return signedToken, nil
}
