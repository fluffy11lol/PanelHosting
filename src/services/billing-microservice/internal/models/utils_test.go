package models

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	secretKey = "testsecret"
)

func TestTokenClaimGeneration(t *testing.T) {
	// Создаем тестовые данные
	claims := TokenClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "testissuer",
		},
		UserID: "12345",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	if signedToken == "" {
		t.Error("Generated token is empty")
	}
}

func TestTokenClaimParsing(t *testing.T) {

	claims := TokenClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "testissuer",
		},
		UserID: "12345",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	parsedToken, err := jwt.ParseWithClaims(signedToken, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(*TokenClaim); ok && parsedToken.Valid {
		if claims.UserID != "12345" {
			t.Errorf("Expected UserID '12345', but got '%s'", claims.UserID)
		}
		if claims.Issuer != "testissuer" {
			t.Errorf("Expected Issuer 'testissuer', but got '%s'", claims.Issuer)
		}
	} else {
		t.Error("Failed to parse claims or token is invalid")
	}
}

func TestTokenClaimExpiration(t *testing.T) {

	claims := TokenClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
			Issuer:    "testissuer",
		},
		UserID: "12345",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	_, err = jwt.ParseWithClaims(signedToken, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err == nil {
		t.Error("Expected error for expired token, but got none")
	}
}
