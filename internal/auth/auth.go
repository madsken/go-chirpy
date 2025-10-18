package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(pw string) (string, error) {
	return argon2id.CreateHash(pw, argon2id.DefaultParams)
}

func CheckPasswordHash(pw, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pw, hash)
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	headerParts := strings.Fields(authHeader)
	if len(headerParts) != 2 || headerParts[0] != "ApiKey" {
		return "", fmt.Errorf("invalid auth header")
	}
	return headerParts[1], nil
}
