// Package passwordutil contains utility functions for working with passwords.
package passwordutil

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"log"
	"strings"
)

// GeneratePasswordHash returns base64 encoded string of the password hash (hashed using SHA512)
func GeneratePasswordHash(password string) (string, error) {

	if len(strings.TrimSpace(password)) <= 0 { // trim the input string and check for empty value
		return "", errors.New("Invalid password")
	}

	passwordByte := []byte(strings.TrimSpace(password))
	hash := sha512.New()
	hash.Write(passwordByte)
	passwordSum := hash.Sum(nil)
	base64Sting := base64.StdEncoding.EncodeToString(passwordSum) // encode to base64 string
	log.Println("Password hash generated successfully")
	return base64Sting, nil
}
