package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/madsken/go-chirpy/internal/database"
)

const port string = "8080"
const fileRootPath string = "."

func initDatabase() (*database.Queries, error) {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return database.New(db), nil
}

func initServerHandlers(apiCfg *apiConfig) *http.ServeMux {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.displayHitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)

	// users endpoint
	serveMux.HandleFunc("POST /api/users", apiCfg.createUser)
	serveMux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginUser)

	// chirps endpoint
	serveMux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)

	// refresh token endpoint
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revoke)

	// polka webhook
	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.polkaWebHook)

	// Create fileserver endpoint
	fs := http.StripPrefix("/app/", http.FileServer(http.Dir(fileRootPath)))
	serveMux.Handle("/app/", apiCfg.mwMetricsInc(fs))

	return serveMux
}

func main() {
	dbQueries, err := initDatabase()
	if err != nil {
		log.Fatalf("Error initalising database: %s", err)
	}

	apiCfg := apiConfig{
		dbQueries: dbQueries,
		platform:  os.Getenv("PLATFORM"),
		secret:    os.Getenv("SECRET"),
		polkaKey:  os.Getenv("POLKA_KEY"),
	}
	serveMux := initServerHandlers(&apiCfg)

	// Construct server
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Starting server on port %s\n", port)

	// Server
	err = server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
		log.Fatal("server had an error")
	}
}
