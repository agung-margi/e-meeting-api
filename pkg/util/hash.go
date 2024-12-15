package util

import (
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/bcrypt"
)

// GenerateSalt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
