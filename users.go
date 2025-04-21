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
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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

	usrRes := userResponse{
		Id:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}

	respondWithJSON(w, 200, usrRes)
}
