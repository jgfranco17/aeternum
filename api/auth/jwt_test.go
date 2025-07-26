package auth

import (
	"os"
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Set up test environment
	os.Setenv("AETERNUM_JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("AETERNUM_JWT_SECRET")

	userID := "test-user-123"
	email := "test@example.com"

	// Generate token
	token, err := GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Validate token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}

	// Check that token is not expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("Token is already expired")
	}
}

func TestValidateInvalidToken(t *testing.T) {
	// Set up test environment
	os.Setenv("AETERNUM_JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("AETERNUM_JWT_SECRET")

	// Test with invalid token
	_, err := ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestGenerateTokenWithoutSecret(t *testing.T) {
	// Ensure no secret is set
	os.Unsetenv("AETERNUM_JWT_SECRET")

	_, err := GenerateToken("user", "email@example.com")
	if err == nil {
		t.Error("Expected error when JWT secret is not set, got nil")
	}
}
