package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"database/sql"
	"errors"

	"github.com/alleviation1/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) upgradeUserHandler(w http.ResponseWriter, r * http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unable to get header key in upgrade user")
		return
	}

	if apiKey != c.polkaAPIKey {
		respondWithError(w, 500, "Invalid api key in upgrade user")
		return
	}

	type requestBody struct {
		Event	string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	req := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 500, "Unable to parse body in upgrade user")
		return
	}

	if req.Event != "user.upgraded" {
		fmt.Printf("Error invalid event in upgrade user: %w\n", err)
		w.WriteHeader(204)
		return
	}

	id, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't parse user id in upgrade user")
		return
	}

	err = c.db.UpgradeUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user in upgrade user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user in upgrade user")
		return
	}

	w.WriteHeader(204)
}