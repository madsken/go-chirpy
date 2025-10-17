package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"
	expiresIn := time.Minute * 5

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error making token: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if parsedID != userID {
		t.Errorf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	userID := uuid.New()
	secret := "correctsecret"
	wrongSecret := "wrongsecret"

	token, err := MakeJWT(userID, secret, time.Minute*5)
	if err != nil {
		t.Fatalf("unexpected error making token: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("expected validation to fail with wrong secret, but got no error")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"

	// Create a token that expired 1 second ago
	token, err := MakeJWT(userID, secret, -1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error making token: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected validation to fail for expired token, but got no error")
	}
}

func TestValidateJWT_InvalidIssuer(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"

	// Create a token with a different issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "evil_issuer",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		Subject:   userID.String(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("unexpected error signing token: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatal("expected validation to fail due to invalid issuer, but got no error")
	}
}

func TestValidateJWT_InvalidTokenFormat(t *testing.T) {
	secret := "supersecret"
	invalidToken := "this.is.not.a.jwt"

	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Fatal("expected validation to fail for invalid token format, but got no error")
	}
}

func TestValidateJWT_InvalidUserID(t *testing.T) {
	secret := "supersecret"

	// Make token with malformed user ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    TokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		Subject:   "not-a-valid-uuid",
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("unexpected error signing token: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatal("expected validation to fail for invalid user ID, but got no error")
	}
}

func TestAuthHeaderStripping(t *testing.T) {
	bearer := "Bearer "
	goodToken := "insane-test-token-420"
	header := http.Header{}
	header.Set("Authorization", bearer+goodToken)

	authToken, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("got error on correct token: %v", err)
	}
	if authToken != goodToken {
		t.Fatalf("unexpected token returned. expected: %v, returned: %v", goodToken, authToken)
	}
}
