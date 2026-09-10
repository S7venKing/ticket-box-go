package auth

import (
	"testing"
	"time"

	jwtpkg "github.com/golang-jwt/jwt/v5"
)

func TestTokenService_GenerateAndValidate(t *testing.T) {
	service := NewTokenService("super-secret-key")

	token, err := service.Generate("user-123", "alice@example.com")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	claims, err := service.Validate(token)
	if err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Fatalf("UserID mismatch: got %q", claims.UserID)
	}

	if claims.Email != "alice@example.com" {
		t.Fatalf("Email mismatch: got %q", claims.Email)
	}
}

func TestTokenService_RejectsExpiredToken(t *testing.T) {
	service := NewTokenService("super-secret-key")

	token, err := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, jwtpkg.MapClaims{
		"sub":   "user-123",
		"email": "alice@example.com",
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
		"exp":   time.Now().Add(-time.Hour).Unix(),
		"type":  "access",
	}).SignedString([]byte("super-secret-key"))
	if err != nil {
		t.Fatalf("SignedString() unexpected error: %v", err)
	}

	if _, err := service.Validate(token); err == nil {
		t.Fatal("Validate() expected error for expired token")
	}
}
