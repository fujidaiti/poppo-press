package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateToken returns a new cryptographically secure random token encoded
// using URL-safe base64 without padding. The raw entropy is 32 bytes.
func GenerateToken() (string, error) {
	const entropyBytes = 32
	b := make([]byte, entropyBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken returns a stable URL-safe base64-encoded SHA-256 digest of the
// provided token. Store only this hash in the database.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
