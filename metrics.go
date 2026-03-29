package main

import (
	"net/http"
	"fmt"
	"strconv"
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