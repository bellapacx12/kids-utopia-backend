package security

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomToken creates a secure random refresh token
func GenerateRandomToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}