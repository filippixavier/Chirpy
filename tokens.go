package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/filippixavier/Chirpy/internal/auth"
	"github.com/filippixavier/Chirpy/internal/database"
)

func (apiCfg *apiConfig) refresh_token(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "token not found", err)
		return
	}

	rtk, err := apiCfg.db.GetRefreshToken(r.Context(), token)

	if err != nil || rtk.RevokedAt.Valid || rtk.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "token not found or revoked", err)
		return
	}

	tk, err := auth.MakeJWT(rtk.UserID, apiCfg.secret, time.Hour)

	if err != nil {
		respondWithError(w, 500, "error when creating token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, response{Token: tk})
}

func (apiCfg *apiConfig) revoke_token(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 400, "no token found", err)
		return
	}

	rtk, err := apiCfg.db.GetRefreshToken(r.Context(), refresh)

	if err != nil {
		respondWithError(w, 400, "token not found", err)
		return
	}

	apiCfg.db.RevokeRefrehToken(r.Context(), database.RevokeRefrehTokenParams{
		Token:     rtk.Token,
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	w.WriteHeader(204)
}
