package main

import (
	"encoding/json"
	"net/http"

	"github.com/filippixavier/Chirpy/internal/auth"
	"github.com/filippixavier/Chirpy/internal/database"
	"github.com/google/uuid"
)

const (
	userUpgraded string = "user.upgraded"
)

func (apiCfg *apiConfig) polkaWebhook(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	var req request

	key, err := auth.GetApiKey(r.Header)

	if err != nil || key != apiCfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error when parsing request body", err)
		return
	}

	defer r.Body.Close()

	if req.Event != userUpgraded {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	usrId, err := uuid.Parse(req.Data.UserId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse user id", err)
		return
	}

	_, err = apiCfg.db.GetUserById(r.Context(), usrId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	_, err = apiCfg.db.UpgradeRedStatus(r.Context(), database.UpgradeRedStatusParams{
		ID:          usrId,
		IsChirpyRed: true,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
