package main

import (
	"net/http"
	"fmt"
	"strconv"
	"os"

	"github.com/joho/godotenv"
)

func (c *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	adminIndex := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %s times!</p>
	</body>
</html>`, strconv.Itoa(int(c.fileserverHits.Load())))
	w.Write([]byte(adminIndex))
}

func (c *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	auth := os.Getenv("PLATFORM")
	if auth != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	err := c.db.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}

	c.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(int(c.fileserverHits.Load()))))
}

func (c *apiConfig) incFileHitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}