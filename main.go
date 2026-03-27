package main

import (
	"net/http"
	"log"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (c *apiConfig) incFileHitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := strconv.Itoa(int(c.fileserverHits.Load()))
	w.Write([]byte("Hits: " + hits))
}

func (c *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(int(c.fileserverHits.Load()))))
}

func main() {
	const port = ":8080"
	const fileRoot = "."
	config := apiConfig{}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", http.StripPrefix("/app/", config.incFileHitsMiddleware(http.FileServer(http.Dir(fileRoot)))))

	serverMux.HandleFunc("GET /healthz", Healthz)
	serverMux.HandleFunc("GET /metrics", config.Metrics)
	serverMux.HandleFunc("POST /reset", config.reset)

	server := &http.Server{
		Addr: port,
		Handler: serverMux,
	}

	log.Fatal(server.ListenAndServe())

	defer server.Close()
}