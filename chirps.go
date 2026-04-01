package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/alleviation1/chirpy/internal/database"
)

func (c *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

	type requestBody struct{
		Body 	    string	  `json:"body"`
		UserID      uuid.UUID `json:"user_id"`
	}

	type responseBody struct{
		ID 			uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Body 		string	  `json:"body"`
		UserID 		uuid.UUID `json:"user_id"`
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

	filteredBody, err := filterBadWords(req.Body)
	if err != nil {
		respondWithError(w, 400, "Unable to censor chirp")
		return
	}
	
	chirp, err := c.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: filteredBody,
		UserID: req.UserID,
	})
	if err != nil {
		respondWithError(w, 500, "Unable to create chirp")
		return
	}


	respondWithValidJson(w, http.StatusCreated, responseBody{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}