package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/madsken/go-chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
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
	if cfg.platform != "dev" {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	writer.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	err := cfg.dbQueries.DeleteAllUsers(request.Context())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error deleting all users", err)
	}
}

// mw to increment hits on /app/ hits
func (cfg *apiConfig) mwMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, request)
	})
}
