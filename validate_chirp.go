package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Chirp struct {
	Body string `json:"body"`
}

type ValidateResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

func (apiCfg *apiConfig) validate_chirp(w http.ResponseWriter, r *http.Request) {
	var chirp Chirp
	response := ValidateResponse{}

	w.Header().Add("Content-Type", "application/json")
	write_json := func() {
		encoded, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(encoded)
	}

	defer r.Body.Close()
	defer write_json()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&chirp); err != nil {
		response.Error = "Something went wrong"
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(chirp.Body) < 140 {
		response.Valid = true
		w.WriteHeader(200)
	} else {
		response.Valid = false
		response.Error = "Chirp is too long"
		w.WriteHeader(400)
	}
}
