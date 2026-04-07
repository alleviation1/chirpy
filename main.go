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
	jwtSecret string
	polkaAPIKey string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	jwtKey := os.Getenv("DB_URL")
	if jwtKey == "" {
		log.Fatal("JWT key must be set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("Polka key must be set")
	}

	const port = ":8080"
	const fileRoot = "."

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	config := apiConfig{
		fileserverHits: atomic.Int32{},
		db: database.New(db),
		jwtSecret: jwtKey,
		polkaAPIKey: polkaKey,
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", config.incFileHitsMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))))

	serverMux.HandleFunc("GET /api/healthz", healthz)
	serverMux.HandleFunc("POST /api/users", config.createUserHandler)
	serverMux.HandleFunc("PUT /api/users", config.updateUserHandler)
	serverMux.HandleFunc("POST /api/refresh", config.refreshTokenHandler)
	serverMux.HandleFunc("POST /api/revoke", config.revokeTokenHandler)
	serverMux.HandleFunc("POST /api/login", config.loginHandler)
	serverMux.HandleFunc("POST /api/chirps", config.createChirpHandler)
	serverMux.HandleFunc("GET /api/chirps", config.getChirpsHandler)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", config.getChirpByIDHandler)
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", config.deleteChirpByIDHandler)
	serverMux.HandleFunc("POST /api/polka/webhooks", config.upgradeUserHandler)

	serverMux.HandleFunc("GET /admin/metrics", config.metrics)
	serverMux.HandleFunc("POST /admin/reset", config.reset)

	server := &http.Server{
		Addr: port,
		Handler: serverMux,
	}

	log.Fatal(server.ListenAndServe())
}