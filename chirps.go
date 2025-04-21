package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/filippixavier/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) chirps(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	var ch chirp

	w.Header().Add("Content-Type", "application/json")

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ch); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(ch.Body) < 140 {
		if _, err := apiCfg.db.GetUserById(r.Context(), ch.UserId); err != nil {
			respondWithError(w, 500, "Unknown user", err)
			return
		}
		chdb, err := apiCfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: purge_bad_words(ch.Body), UserID: ch.UserId})
		if err != nil {
			respondWithError(w, 500, "Error when inserting chirp in db", err)
		}
		res := response{
			Id:        chdb.ID,
			CreatedAt: chdb.CreatedAt,
			UpdatedAt: chdb.UpdatedAt,
			Body:      chdb.Body,
			UserId:    ch.UserId,
		}
		respondWithJSON(w, 201, res)
	} else {
		respondWithError(w, 400, "Chirp is too long", nil)
	}
}

func purge_bad_words(str string) string {
	reg, err := regexp.Compile(`(?i)(kerfuffle|sharbert|fornax)\s*?`)

	if err != nil {
		return str
	}

	res := reg.ReplaceAll([]byte(str), []byte("****"))

	return string(res)
}
