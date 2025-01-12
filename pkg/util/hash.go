package util

import (
	"crypto/rand"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := strings.Builder{}
	for i := 0; i < length; i++ {
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result.WriteByte(charset[randomInt.Int64()])
	}
	return result.String()
}
