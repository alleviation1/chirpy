package main

import (
	"net/http"
	"encoding/json"
	"time"
	"fmt"

	"github.com/alleviation1/chirpy/internal/auth"
	"github.com/alleviation1/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email 	  string `json:"email"`
		Password  string `json:"password"`
	}

	type responseBody struct {
		ID 		  uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email 	  string    `json:"email"`
		IsChirpyRed	bool	`json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 500, "Unable to process request in create user")
		return
	}

	hashedPass, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password in create user")
		return
	}

	user, err := c.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, 500, "Unable to create user in create user")
		return
	}

	resp := responseBody{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithValidJson(w, 201, resp)
}

func (c *apiConfig) updateUserHandler (w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email 	 string `json:"email"`
		Password string `json:"password"`
	}

	type UpdatedUser struct {
		ID			uuid.UUID    `json:"id"`
		Email		string       `json:"email"`
		CreatedAt	time.Time    `json:"created_at"`
		UpdatedAt	time.Time    `json:"updated_at"`
		IsChirpyRed	bool		 `json:"is_chirpy_red"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unable to get auth header in update user handler")
		return
	}

	req := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	defer r.Body.Close()

	userID, err := auth.ValidateJWT(tokenString, c.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unable to validate user jwt in update user handler")
		return
	}

	hashedPass, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 500, "Unable to hash password in update user handler")
		return
	}

	updatedUser, err := c.db.SetEmailAndPassword(r.Context(), database.SetEmailAndPasswordParams {
		Email: req.Email,
		HashedPassword: hashedPass,
		ID: userID,
	})
	if err != nil {
		fmt.Printf("Error: %w\n", err)
		respondWithError(w, 500, "Unable to set new email and password in update user handler")
		return
	}

	respondWithValidJson(w, 200, UpdatedUser{
		ID: updatedUser.ID,
		Email: updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}