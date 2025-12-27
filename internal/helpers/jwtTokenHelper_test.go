package helpers

import (
	"testing"

	"main.go/internal/models"
)

func TestCreateAndVerifyToken(t *testing.T) {
	SecretKey = []byte("test-secret")

	user := &models.User{
		Id:   123,
		Name: "someone",
	}

	token, err := CreateToken(user)
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}
	if claims.Subject != "123" {
		t.Fatalf("expected subject 123, got %q", claims.Subject)
	}
	if claims.Username != "someone" {
		t.Fatalf("expected username someone, got %q", claims.Username)
	}
}

func TestVerifyTokenInvalid(t *testing.T) {
	SecretKey = []byte("test-secret")

	if _, err := VerifyToken("not-a-token"); err == nil {
		t.Fatal("expected error for invalid token")
	}
}
