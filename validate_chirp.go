package main

import (
	"net/http"
	"encoding/json"
	"strings"
)

func validate_chirp_handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct{
		Body string	`json:"body"`
	}

	type responseBody struct{
		Valid bool	 		`json:"valid"`
		Cleaned_Body string	`json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 500, "Unable to decode JSON in chirp handler")
		return
	}

	if len(string(req.Body)) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	filteredResp, err := filterBadWords(req.Body)
	if err != nil {
		respondWithError(w, 400, "Unable to censor chirp")
		return
	}
	

	resp := responseBody{
		Valid: true,
		Cleaned_Body: filteredResp,
	}
	respondWithValidJson(w, 200, resp)
}

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