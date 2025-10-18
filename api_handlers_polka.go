package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/madsken/go-chirpy/internal/auth"
	"github.com/madsken/go-chirpy/internal/database"
)

func validateAPIKey(apiKey string, header http.Header) error {
	headerKey, err := auth.GetAPIKey(header)
	if err != nil {
		return err
	}

	if headerKey != apiKey {
		return fmt.Errorf("header api key does not match env key")
	}
	return nil
}

func (cfg *apiConfig) polkaWebHook(writer http.ResponseWriter, request *http.Request) {
	err := validateAPIKey(cfg.polkaKey, request.Header)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "invalid header key", err)
		return
	}

	type reqJson struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	reqData := reqJson{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&reqData)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	if reqData.Event != "user.upgraded" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	updatedUser, err := cfg.dbQueries.UpgradeUserChirpyRedStatus(request.Context(), database.UpgradeUserChirpyRedStatusParams{
		IsChirpyRed: true,
		ID:          reqData.Data.UserID,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error upgrading status", err)
		return
	}

	if updatedUser.ID != reqData.Data.UserID {
		respondWithError(writer, http.StatusNotFound, "User not found", nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
