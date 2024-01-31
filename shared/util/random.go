package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

func RandomSalt() string {
	salt := make([]byte, 16)

	if os.Getenv("TEST") == "TRUE" {
		return "saltedtest"
	}

	_, err := rand.Read(salt)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(salt)
}

func HashPassword(s string) string {
	if os.Getenv("TEST") == "TRUE" {
		return "hashedpassworduser123123123"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(hashedPassword)
}

func IsPasswordCorrect(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func IsCorrectPhoneNumber(phoneNumber string) bool {
	return !strings.HasPrefix(phoneNumber, "+62")
}

func ValidatePhoneNumber(phoneNumber string) error {
	if !strings.HasPrefix(phoneNumber, "+62") {
		return errors.New("phone number should use Indonesia Country Code +62")
	}
	return nil
}
