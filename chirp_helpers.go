package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/madsken/go-chirpy/internal/auth"
)

var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func cleanProfanity(body string) string {
	words := strings.Fields(body)

	badWordsMap := map[string]bool{}
	for _, badWord := range badWords {
		badWordsMap[badWord] = true
	}

	for i, word := range words {
		if badWordsMap[strings.ToLower(word)] {
			words[i] = "****"
		}

	}

	return strings.Join(words, " ")
}

func validateToken(request *http.Request, secretToken string) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		return uuid.Nil, err
	}
	return auth.ValidateJWT(token, secretToken)
}
