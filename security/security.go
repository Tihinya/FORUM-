package security

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const (
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()-=_+[]{}|;:,.<>/?"
)

func PasswordEncrypting(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CreateRandomPassword(length int) string {
	charset := uppercaseLetters + lowercaseLetters + digits + specialChars
	password := make([]byte, length)

	for i := range password {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password)
}
