package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/filippixavier/Chirpy/internal/auth"
	"github.com/filippixavier/Chirpy/internal/database"
	"github.com/google/uuid"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (apiCfg *apiConfig) create_user(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var usrJson userRequest
	ctx := r.Context()

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&usrJson); err != nil {
		respondWithError(w, 500, "error when decoding json", err)
		return
	}

	hpwd, err := auth.HashPassword(usrJson.Password)

	if err != nil {
		respondWithError(w, 500, "error when hashing password", err)
		return
	}

	usr, err := apiCfg.db.CreateUser(ctx, database.CreateUserParams{Email: usrJson.Email, HashedPassword: hpwd})

	if err != nil {
		respondWithError(w, 500, "error when creating user", err)
		return
	}

	usrRes := userResponse{
		Id:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}

	respondWithJSON(w, 201, usrRes)
}

func (apiCfg *apiConfig) login_user(w http.ResponseWriter, r *http.Request) {
	var usrJson userRequest
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&usrJson); err != nil {
		respondWithError(w, 500, "error when decoding json", err)
		return
	}

	usr, err := apiCfg.db.GetUserByEmail(r.Context(), usrJson.Email)

	if err != nil {
		respondWithError(w, 401, "incorrect username or password", err)
		return
	}

	err = auth.CheckPasswordHash(usr.HashedPassword, usrJson.Password)

	if err != nil {
		respondWithError(w, 401, "incorrect username or password", err)
		return
	}

	token, err := auth.MakeJWT(usr.ID, apiCfg.secret, time.Hour)

	if err != nil {
		respondWithError(w, 500, "error when creating token", err)
		return
	}

	refresh, _ := auth.MakeRefreshToken()

	rtk, err := apiCfg.db.MakeRefreshToken(r.Context(), database.MakeRefreshTokenParams{
		Token:     refresh,
		UserID:    usr.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})

	if err != nil {
		respondWithError(w, 500, "error when creating refresh token", err)
		return
	}

	usrRes := userResponse{
		Id:           usr.ID,
		CreatedAt:    usr.CreatedAt,
		UpdatedAt:    usr.UpdatedAt,
		Email:        usr.Email,
		Token:        token,
		RefreshToken: rtk.Token,
	}

	respondWithJSON(w, 200, usrRes)
}

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
