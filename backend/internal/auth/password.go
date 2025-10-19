// Package auth contains helpers related to authentication concerns such as
// password hashing and verification.
package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// HashPassword derives a password hash using argon2id and returns a
// PHC-style encoded string containing parameters, salt, and hash.
func HashPassword(password string) (string, error) {
	const (
		time    = 1
		memory  = 64 * 1024
		threads = 4
		keyLen  = 32
		saltLen = 16
	)

	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	// PHC string: $argon2id$v=19$m=65536,t=1,p=4$<salt_b64>$<hash_b64>
	phc := "$argon2id$v=19$m=65536,t=1,p=4$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
	return phc, nil
}
