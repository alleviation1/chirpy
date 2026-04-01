package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// get req body
	defer r.Body.Close()

	type requestBody struct {
		Email 	  string `json:"email"`
	}

	type responseBody struct {
		ID 		  uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email 	  string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 500, "Unable to process request in create user")
		return
	}

	user, err := c.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, 500, "Unable to create user in create user")
		return
	}

	resp := responseBody{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}


	respondWithValidJson(w, 201, resp)
}