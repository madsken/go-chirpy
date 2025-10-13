package main

import (
	"strings"
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
