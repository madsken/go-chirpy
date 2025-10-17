package main

import (
	"encoding/json"
	"net/http"

	"github.com/madsken/go-chirpy/internal/auth"
	"github.com/madsken/go-chirpy/internal/database"
)

func (cfg *apiConfig) createUser(writer http.ResponseWriter, request *http.Request) {
	type reqJson struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	reqData := reqJson{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	hashedPw, err := auth.HashPassword(reqData.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(request.Context(), database.CreateUserParams{
		Email:    reqData.Email,
		Password: hashedPw,
	})
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

func (cfg *apiConfig) loginUser(writer http.ResponseWriter, request *http.Request) {
	type reqJson struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	reqData := reqJson{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(request.Context(), reqData.Email)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error getting user from database", err)
		return
	}

	ok, err := auth.CheckPasswordHash(reqData.Password, user.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error checking hash and pw", err)
		return
	}
	if !ok {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(writer, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
