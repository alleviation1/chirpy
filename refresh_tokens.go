package main

import (
	"net/http"

	"github.com/alleviation1/chirpy/internal/auth"
)

func (c *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	type respBody struct{
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Unable to get auth header in refresh token handler")
		return
	}

	userID, err := c.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, 401, "Unable to find user or invalid token in refresh token handler")
		return
	}

	newToken, err := auth.MakeJWT(userID, c.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unable to make new access token in refresh token handler")
		return
	}

	respondWithValidJson(w, 200, respBody{
		Token: newToken,
	})
}

func (c *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "Unable to get auth header in revoke token handler")
		return
	}

	err = c.db.RevokeToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, 500, "Unable to revoke token in revoke token handler")
		return
	}

	respondWithValidJson(w, 204, "")
}