package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(data any) (*string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(fmt.Appendf(nil, "%v", data))
	if err != nil {
		return nil, fmt.Errorf("failed to hash token: %w", err)
	}
	token := hex.EncodeToString(hasher.Sum(nil))
	return &token, nil
}
func GenerateHashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash data: %w", err)
	}

	hashed := string(hash)
	return &hashed, nil
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
