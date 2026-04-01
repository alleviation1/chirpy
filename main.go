package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/alleviation1/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	const port = ":8080"
	const fileRoot = "."

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error creating postgres connection: %w", err)
	}
	
	defer db.Close()

	config := apiConfig{
		fileserverHits: atomic.Int32{},
		db: database.New(db),
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", config.incFileHitsMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))

	serverMux.HandleFunc("GET /api/healthz", healthz)
	serverMux.HandleFunc("POST /api/users", config.createUserHandler)
	serverMux.HandleFunc("POST /api/chirps", config.createChirpHandler)
	serverMux.HandleFunc("GET /api/chirps", config.getChirpsHandler)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", config.getChirpByIDHandler)

	serverMux.HandleFunc("GET /admin/metrics", config.metrics)
	serverMux.HandleFunc("POST /admin/reset", config.reset)

	server := &http.Server{
		Addr: port,
		Handler: serverMux,
	}

	log.Fatal(server.ListenAndServe())
}