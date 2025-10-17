package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	errResponse := errorResponse{
		Error: msg,
	}

	respondWithJSON(writer, statusCode, errResponse)
}

func respondWithJSON(writer http.ResponseWriter, statusCode int, payload interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling response to JSON: %s", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(statusCode)
	writer.Write(data)
}
