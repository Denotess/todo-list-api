package helpers

import "testing"

func TestHashPasswordAndCheck(t *testing.T) {
	password := "test-password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == password {
		t.Fatal("hash should not match the raw password")
	}

	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected password to match hash")
	}

	ok, err = CheckPasswordHash("wrong-password", hash)
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	if ok {
		t.Fatal("expected password mismatch")
	}
}
