package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the plain password
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) // or handle error properly
	}
	return string(hash)
}

// CheckPassword compares plain password with hashed password
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
