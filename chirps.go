package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/filippixavier/Chirpy/internal/auth"
	"github.com/filippixavier/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (apiCfg *apiConfig) create_chirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	var ch chirp

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unauthorized access", err)
		return
	}

	usr, err := auth.ValidateJWT(token, apiCfg.secret)

	if err != nil {
		respondWithError(w, 401, "unauthorized access", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ch); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(ch.Body) < 140 {
		if _, err := apiCfg.db.GetUserById(r.Context(), usr); err != nil {
			respondWithError(w, 500, "Unknown user", err)
			return
		}
		chdb, err := apiCfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: purge_bad_words(ch.Body), UserID: usr})
		if err != nil {
			respondWithError(w, 500, "Error when inserting chirp in db", err)
			return
		}
		res := response{
			Id:        chdb.ID,
			CreatedAt: chdb.CreatedAt,
			UpdatedAt: chdb.UpdatedAt,
			Body:      chdb.Body,
			UserId:    usr,
		}
		respondWithJSON(w, 201, res)
	} else {
		respondWithError(w, 400, "Chirp is too long", nil)
	}
}

func (apiCfg *apiConfig) delete_chirp(w http.ResponseWriter, r *http.Request) {
	tkn, err := auth.GetBearerToken(r.Header)
	chipPath := r.PathValue("chirpID")

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token found", err)
		return
	}

	usrId, err := auth.ValidateJWT(tkn, apiCfg.secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirpID, err := uuid.Parse(chipPath)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}
	chirp, err := apiCfg.db.GetChirpById(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if chirp.UserID != usrId {
		respondWithError(w, http.StatusForbidden, "not the owner of chirp", err)
		return
	}

	if _, err := apiCfg.db.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error when deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (apiCfg *apiConfig) get_chirps(w http.ResponseWriter, r *http.Request) {
	type response []Chirp

	var chirps []database.Chirp
	var err error

	usr := r.URL.Query().Get("author_id")
	ord := r.URL.Query().Get("sort")

	if usr == "" {
		chirps, err = apiCfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, "Unknown error", err)
			return
		}
	} else {
		usrId, err := uuid.Parse(usr)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "invalid author id", err)
			return
		}
		chirps, err = apiCfg.db.GetChirpsByAuthor(r.Context(), usrId)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error when fetching result", err)
			return
		}
	}

	chJson := make(response, len(chirps))

	for i, ch := range chirps {
		chJson[i].ID = ch.ID
		chJson[i].CreatedAt = ch.CreatedAt
		chJson[i].UpdatedAt = ch.UpdatedAt
		chJson[i].Body = ch.Body
		chJson[i].UserID = ch.UserID
	}

	if ord == "desc" {
		slices.SortFunc(chJson, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}

	respondWithJSON(w, 200, chJson)
}

func (apiCfg *apiConfig) get_chirp_by_id(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")

	uid, err := uuid.Parse(id)

	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	chirp, err := apiCfg.db.GetChirpById(r.Context(), uid)

	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func purge_bad_words(str string) string {
	reg, err := regexp.Compile(`(?i)(kerfuffle|sharbert|fornax)\s*?`)

	if err != nil {
		return str
	}

	res := reg.ReplaceAll([]byte(str), []byte("****"))

	return string(res)
}
