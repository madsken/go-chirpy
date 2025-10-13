package main

import (
	"fmt"
	"log"
	"net/http"
)

const port string = "8080"
const fileRootPath string = "."

func main() {
	serveMux := http.NewServeMux()
	apiCfg := apiConfig{}

	// Create readiness endpoint
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.displayHitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirp)

	// Create fileserver endpoint
	fs := http.StripPrefix("/app/", http.FileServer(http.Dir(fileRootPath)))
	serveMux.Handle("/app/", apiCfg.mwMetricsInc(fs))

	// Construct server
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Starting server on port %s\n", port)

	// Server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
		log.Fatal("server had an error")
	}
}
