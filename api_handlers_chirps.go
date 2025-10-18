package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/madsken/go-chirpy/internal/auth"
	"github.com/madsken/go-chirpy/internal/database"
)

func (cfg *apiConfig) createChirp(writer http.ResponseWriter, request *http.Request) {
	type reqJson struct {
		Body string `json:"body"`
	}

	userID, err := validatePost(request, cfg.secret)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "invalid token", err)
		return
	}

	reqData := reqJson{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(request.Context(), database.CreateChirpParams{
		Body:   reqData.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	// Check chirp lenght
	if len(chirp.Body) > 140 {
		respondWithError(writer, http.StatusBadRequest, "Chirp too long", nil)
		return
	}
	chirp.Body = cleanProfanity(chirp.Body)

	respondWithJSON(writer, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreateAt:  chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) getChirps(writer http.ResponseWriter, request *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(request.Context())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error fetching chirps", err)
		return
	}

	chirpsResp := []Chirp{}
	for _, chirp := range chirps {
		chirpsResp = append(chirpsResp, Chirp{
			ID:        chirp.ID,
			CreateAt:  chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(writer, http.StatusOK, chirpsResp)
}

func (cfg *apiConfig) getChirp(writer http.ResponseWriter, request *http.Request) {
	cID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(request.Context(), cID)
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(writer, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreateAt:  chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validatePost(request *http.Request, secretToken string) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		return uuid.Nil, nil
	}
	return auth.ValidateJWT(token, secretToken)
}
