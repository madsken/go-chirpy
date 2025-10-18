package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/madsken/go-chirpy/internal/database"
)

func (cfg *apiConfig) createChirp(writer http.ResponseWriter, request *http.Request) {
	type reqJson struct {
		Body string `json:"body"`
	}

	userID, err := validateToken(request, cfg.secret)
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
	authorID, _ := uuid.Parse(request.URL.Query().Get("author_id"))
	sortBy := request.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "asc"
	}

	chirps, err := cfg.dbQueries.GetAllChirps(request.Context())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error fetching chirps", err)
		return
	}

	chirpsResp := []Chirp{}
	for _, chirp := range chirps {
		if chirp.UserID != authorID && authorID != uuid.Nil {
			continue
		}
		chirpsResp = append(chirpsResp, Chirp{
			ID:        chirp.ID,
			CreateAt:  chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	sort.Slice(chirpsResp, func(i, j int) bool {
		if sortBy == "desc" {
			return chirpsResp[i].CreateAt.After(chirpsResp[j].CreateAt)
		}
		return chirpsResp[i].CreateAt.Before(chirpsResp[j].CreateAt)
	})

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

func (cfg *apiConfig) deleteChirp(writer http.ResponseWriter, request *http.Request) {
	userID, err := validateToken(request, cfg.secret)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "invalid token", err)
		return
	}
	cID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}

	chirpData, err := cfg.dbQueries.GetChirp(request.Context(), cID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(writer, http.StatusNotFound, "chirp not in database", err)
			return
		}
		respondWithError(writer, http.StatusInternalServerError, "error getting chirp from database", err)
		return
	}

	if chirpData.UserID != userID {
		respondWithError(writer, http.StatusForbidden, "user does not own chirp", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(request.Context(), cID)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "failed to delete chirp", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
