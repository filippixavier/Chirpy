package main

import (
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
	IsChirpyRed  bool      `json:"is_chirpy_red"`
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

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&usrJson); err != nil {
		respondWithError(w, 500, "error when decoding json", err)
		return
	}
	defer r.Body.Close()

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
		Id:          usr.ID,
		Email:       usr.Email,
		IsChirpyRed: usr.IsChirpyRed,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
	}

	respondWithJSON(w, 201, usrRes)
}

func (apiCfg *apiConfig) login_user(w http.ResponseWriter, r *http.Request) {
	var usrJson userRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&usrJson); err != nil {
		respondWithError(w, 500, "error when decoding json", err)
		return
	}
	defer r.Body.Close()

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
		Email:        usr.Email,
		IsChirpyRed:  usr.IsChirpyRed,
		CreatedAt:    usr.CreatedAt,
		UpdatedAt:    usr.UpdatedAt,
		Token:        token,
		RefreshToken: rtk.Token,
	}

	respondWithJSON(w, 200, usrRes)
}

func (apiCfg *apiConfig) update_user(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var body request

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "missing auth token", err)
		return
	}

	id, err := auth.ValidateJWT(token, apiCfg.secret)

	if err != nil {
		respondWithError(w, 401, "invalid auth token", err)
		return
	}

	decode := json.NewDecoder(r.Body)

	err = decode.Decode(&body)

	if err != nil {
		respondWithError(w, 500, "error when reading body", err)
		return
	}

	defer r.Body.Close()

	hpassword, err := auth.HashPassword(body.Password)

	if err != nil {
		respondWithError(w, 500, "error hashing password", err)
		return
	}

	usr, err := apiCfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             id,
		Email:          body.Email,
		HashedPassword: hpassword,
	})

	if err != nil {
		respondWithError(w, 500, "error when updating user", err)
		return
	}

	respondWithJSON(w, 200, userResponse{
		Id:          usr.ID,
		IsChirpyRed: usr.IsChirpyRed,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
		Email:       usr.Email,
	})
}
