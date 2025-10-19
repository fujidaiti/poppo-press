package auth_test

import (
	"testing"

	"github.com/fujidaiti/poppo-press/backend/internal/auth"
)

func TestHashAndVerifyPassword(t *testing.T) {
	phc, err := auth.HashPassword("secret-123")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	ok, err := auth.VerifyPassword("secret-123", phc)
	if err != nil || !ok {
		t.Fatalf("verify expected ok; ok=%v err=%v", ok, err)
	}
	ok, err = auth.VerifyPassword("wrong", phc)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if ok {
		t.Fatalf("verify expected false for wrong password")
	}
}
