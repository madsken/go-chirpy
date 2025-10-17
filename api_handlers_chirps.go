package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/madsken/go-chirpy/internal/database"
)

func (cfg *apiConfig) createChirp(writer http.ResponseWriter, request *http.Request) {
	type reqJson struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	reqData := reqJson{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(request.Context(), database.CreateChirpParams{
		Body:   reqData.Body,
		UserID: reqData.UserID,
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

	chirpResp := Chirp{
		ID:        chirp.ID,
		CreateAt:  chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(writer, http.StatusCreated, chirpResp)
}

func (cfg *apiConfig) getChirps(writer http.ResponseWriter, request *http.Request) {

}
