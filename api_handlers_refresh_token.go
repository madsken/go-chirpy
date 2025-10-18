package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/madsken/go-chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(writer http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "error getting token from header", err)
		return
	}

	tokenDb, err := cfg.dbQueries.CheckValidToken(request.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(writer, http.StatusUnauthorized, "token not in database", err)
			return
		}
		respondWithError(writer, http.StatusInternalServerError, "error in check valid token query", err)
		return
	}

	newToken, err := auth.MakeJWT(tokenDb.UserID, cfg.secret)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error creating new token", err)
		return
	}

	type NewTokenResp struct {
		Token string `json:"token"`
	}
	respondWithJSON(writer, http.StatusOK, NewTokenResp{
		Token: newToken,
	})

}

func (cfg *apiConfig) revoke(writer http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "error getting token from header", err)
		return
	}

	err = cfg.dbQueries.RevokeToken(request.Context(), token)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "error deleting token in database", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
