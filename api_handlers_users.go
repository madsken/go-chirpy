package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) createUser(writer http.ResponseWriter, request *http.Request) {
	type reqJson struct {
		Email string `json:"email"`
	}
	reqData := reqJson{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(request.Context(), reqData.Email)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error creating user in database", err)
		return
	}

	response := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(writer, http.StatusCreated, response)
}
