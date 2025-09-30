package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a JWT for a given user ID, signed with the mock secret.
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": float64(userID), // JWT numbers are float64 by default
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
		"iat":     time.Now().Unix(), // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte("mock-secret") // Must match middleware secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}