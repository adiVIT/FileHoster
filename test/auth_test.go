package test

import (
	"testing"
	"filestore/internal/auth"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !auth.CheckPasswordHash(password, hashedPassword) {
		t.Errorf("Expected password to match")
	}
}

func TestGenerateJWT(t *testing.T) {
	userID := uint(1)
	token, err := auth.GenerateJWT(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	parsedUserID, err := auth.ValidateJWT(token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("Expected userID %d, got %d", userID, parsedUserID)
	}
}