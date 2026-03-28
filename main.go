package main

import (
	"net/http"
	"log"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = ":8080"
	const fileRoot = "."
	config := apiConfig{fileserverHits: atomic.Int32{}}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", config.incFileHitsMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))

	serverMux.HandleFunc("GET /api/healthz", healthz)

	serverMux.HandleFunc("GET /admin/metrics", config.metrics)
	serverMux.HandleFunc("POST /admin/reset", config.reset)

	server := &http.Server{
		Addr: port,
		Handler: serverMux,
	}

	log.Fatal(server.ListenAndServe())

	defer server.Close()
}