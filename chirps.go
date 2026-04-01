package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/alleviation1/chirpy/internal/database"
)

type Chirp struct{
	ID 			uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Body 		string	  `json:"body"`
	UserID 		uuid.UUID `json:"user_id"`
}

func (c *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct{
		Body 	    string	  `json:"body"`
		UserID      uuid.UUID `json:"user_id"`
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


	respondWithValidJson(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func (c *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chirps, err := c.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get chirps")
		return
	}

	parsedChirps := []Chirp{}
	for _, chirp := range chirps {
		parsedChirps = append(parsedChirps, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithValidJson(w, http.StatusOK, parsedChirps)
}

func (c *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
		type CreateChirpParams struct {
		Body   string
		UserID uuid.UUID
	}
	defer r.Body.Close()

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusInternalServerError, "Chirp ID was not passed correctly")
		return
	}

	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse chirp id into uuid")
		return
	}

	chirp, err := c.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Could not get chirp in get chirp by id")
		return
	}

	respondWithValidJson(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}