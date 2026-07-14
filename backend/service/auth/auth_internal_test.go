package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "MySecretPassword123!"

	hashed, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}
	if hashed == password {
		t.Fatalf("Expected hashed password to differ from plain text")
	}

	if !checkPassword(password, hashed) {
		t.Fatal("Expected password verification to succeed")
	}

	if checkPassword("wrong-password", hashed) {
		t.Fatal("Expected password verification to fail for wrong password")
	}
}

func TestCreateJWT(t *testing.T) {
	tokenString, err := createJWT()
	if err != nil {
		t.Fatalf("createJWT() failed: %v", err)
	}
	if tokenString == "" {
		t.Fatal("Expected non-empty token string")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return getJWTKey(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
	if !token.Valid {
		t.Fatal("Expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Expected MapClaims")
	}

	if _, ok := claims["jti"]; !ok {
		t.Fatal("Expected jti claim")
	}
	if _, ok := claims["nbf"]; !ok {
		t.Fatal("Expected nbf claim")
	}
	if _, ok := claims["exp"]; !ok {
		t.Fatal("Expected exp claim")
	}
}
