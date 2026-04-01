package main

import (
	"net/http"
	"strings"
	"encoding/json"
)

func respondWithValidJson(w http.ResponseWriter, code int, payload interface{}) error {
	dat, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithValidJson(w, code, map[string]string{"error": "Something went wrong"})
}

func filterBadWords(payload string) (string, error) {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(payload, " ")
	for _, badWord := range badWords {
		for j, word := range words {
			if strings.ToLower(word) == badWord {
				words[j] = "****"
			}
		}
	}

	return strings.Join(words, " "), nil
}