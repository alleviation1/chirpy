package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/alleviation1/chirpy/internal/auth"
	"github.com/alleviation1/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) loginHandler (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email 	  string `json:"email"`
		Password  string `json:"password"`
	}

	type User struct {
		ID 		  uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email 	  string `json:"email"`
	}

	type responseBody struct {
		User
		Token 	  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	req := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse request body in login handler")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, 401, "Cannot use empty fields in login handler")
		return
	}

	user, err := c.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get user by email in login handler")
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to check password and hash in login handler")
		return
	}

	if match == false {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, c.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to make jwt in login")
		return
	}

	refreshToken, err := c.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: auth.MakeRefreshToken(),
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Unable to create refresh token")
		return
	}

	respondWithValidJson(w, 200, responseBody{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
		Token: accessToken,
		RefreshToken: refreshToken.Token,
	})
}