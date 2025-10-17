package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(pw string) (string, error) {
	return argon2id.CreateHash(pw, argon2id.DefaultParams)
}

func CheckPasswordHash(pw, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pw, hash)
}
