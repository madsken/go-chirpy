package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) displayHitsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	content := fmt.Sprintf(`
<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, cfg.fileserverHits.Load())

	writer.Write([]byte(content))
}

func (cfg *apiConfig) resetHitsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)
}

// mw to increment hits on /app/ hits
func (cfg *apiConfig) mwMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, request)
	})
}

func validateChirp(writer http.ResponseWriter, request *http.Request) {
	type chirpData struct {
		Body string `json:"body"`
	}
	chirp := chirpData{}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	// Check chirp lenght
	if len(chirp.Body) > 140 {
		respondWithError(writer, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	cleanedBody := cleanProfanity(chirp.Body)

	type responseStruct struct {
		CleanedBody string `json:"cleaned_body"`
	}
	response := responseStruct{
		CleanedBody: cleanedBody,
	}

	respondWithJSON(writer, http.StatusOK, response)
}
