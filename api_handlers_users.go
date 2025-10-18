package main

import (
	"encoding/json"
	"net/http"
	"time"

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

func (cfg *apiConfig) updateUser(writer http.ResponseWriter, request *http.Request) {
	userID, err := validateToken(request, cfg.secret)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "invalid token", err)
		return
	}

	type reqJson struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	reqData := reqJson{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	hashedPw, err := auth.HashPassword(reqData.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updatedRow, err := cfg.dbQueries.UpdateEmailAndPassword(request.Context(), database.UpdateEmailAndPasswordParams{
		Email:    reqData.Email,
		Password: hashedPw,
		ID:       userID,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error updating database with new password and email", err)
		return
	}

	respondWithJSON(writer, http.StatusOK, User{
		ID:        updatedRow.ID,
		CreatedAt: updatedRow.CreatedAt,
		UpdatedAt: updatedRow.UpdatedAt,
		Email:     updatedRow.Email,
	})
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

	token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error creating token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}
	_, err = cfg.dbQueries.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error creating refresh token in database", err)
		return
	}

	respondWithJSON(writer, http.StatusOK, UserLogin{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
