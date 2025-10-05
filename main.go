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

	// Create readiness endpoint
	serveMux.HandleFunc("/healthz", readinessHandler)

	// Create fileserver endpoint
	fs := http.FileServer(http.Dir(fileRootPath))
	serveMux.Handle("/app/", http.StripPrefix("/app/", fs))

	// Construct server
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Starting server on port %s\n", port)

	// Server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("server had an error")
	}
}

func readinessHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK\n"))
}
