package main

import (
	"crypto/sha256"
	"fmt"
)

type User struct {
	Username     string
	PasswordHash string
	Email        string
}

func hashPassword(password string) string {
	// Use a strong hashing algorithm like SHA-256
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
