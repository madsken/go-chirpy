package main

import (
	"testing"
)

func TestChirpProfanity(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "I had something interesting for breakfast",
			expected: "I had something interesting for breakfast",
		},
		{
			input:    "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			input:    "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
	}

	for _, c := range cases {
		actual := cleanProfanity(c.input)

		if actual != c.expected {
			t.Errorf("Actual result not macthing expected: got %v : expected %v", actual, c.expected)
		}
	}
}
