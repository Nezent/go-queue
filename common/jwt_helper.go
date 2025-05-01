package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("nezent-go-queue") // Replace with env in real use

// GenerateJWT generates a JWT with user ID and role
func GenerateJWT(userID string, role string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT parses the token and returns userID and role
func ParseJWT(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("user_id missing in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("role missing in token")
	}

	return userID, role, nil
}
