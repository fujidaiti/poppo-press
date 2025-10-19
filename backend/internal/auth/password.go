// Package auth contains helpers related to authentication concerns such as
// password hashing and verification.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

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

// VerifyPassword checks whether the provided password matches the given PHC
// formatted argon2id hash. It returns true on match and false otherwise.
func VerifyPassword(password, phc string) (bool, error) {
	// Expected PHC format: $argon2id$v=19$m=<mem>,t=<time>,p=<threads>$<salt_b64>$<hash_b64>
	parts := strings.Split(phc, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid phc format")
	}
	// parts[2] == v=19 (ignored)
	params := parts[3]
	var timeCost, memoryKB, threads int
	// defaults in case parsing fails
	timeCost = 1
	memoryKB = 64 * 1024
	threads = 4

	for _, kv := range strings.Split(params, ",") {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			continue
		}
		switch pair[0] {
		case "t":
			if v, err := strconv.Atoi(pair[1]); err == nil {
				timeCost = v
			}
		case "m":
			if v, err := strconv.Atoi(pair[1]); err == nil {
				memoryKB = v
			}
		case "p":
			if v, err := strconv.Atoi(pair[1]); err == nil {
				threads = v
			}
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	got := argon2.IDKey([]byte(password), salt, uint32(timeCost), uint32(memoryKB), uint8(threads), uint32(len(want)))
	if len(got) != len(want) {
		return false, nil
	}
	// constant-time compare
	var diff byte
	for i := range got {
		diff |= got[i] ^ want[i]
	}
	return diff == 0, nil
}
