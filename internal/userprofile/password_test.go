package userprofile_test

import (
	"testing"

	"github.com/bobfive1/user-management-api/internal/userprofile"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "secret-password"

	hash := userprofile.HashPassword(password)
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == password {
		t.Fatal("hash should not equal plain password")
	}
	if !userprofile.CheckPasswordHash(password, hash) {
		t.Fatal("password should match generated hash")
	}
	if userprofile.CheckPasswordHash("wrong-password", hash) {
		t.Fatal("wrong password should not match hash")
	}
}
