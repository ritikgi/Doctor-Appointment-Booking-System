package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID uint, email, role, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
