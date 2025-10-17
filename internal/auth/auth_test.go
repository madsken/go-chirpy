package auth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	cases := []struct {
		pw string
	}{
		{
			pw: "password",
		},
		{
			pw: "reallylongpasswordthatissoveryverylong",
		},
		{
			pw: "f23hjfAOWIFpasd982",
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.pw)
		if err != nil {
			t.Errorf("Error hashing password: %s", err)
		}

		ok, err := CheckPasswordHash(c.pw, hash)
		if err != nil {
			t.Errorf("Error hashing check password: %s", err)
		}
		if !ok {
			t.Errorf("Password check failed")
		}
	}
}
