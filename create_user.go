package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) create_user(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email string `json:"email"`
	}

	type userResponse struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")

	var usrJson userRequest
	ctx := r.Context()

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&usrJson); err != nil {
		respondWithError(w, 500, "error when decoding json", err)
		return
	}

	usr, err := apiCfg.db.CreateUser(ctx, usrJson.Email)

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
