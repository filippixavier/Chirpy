package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func (apiCfg *apiConfig) validate_chirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type validateResponse struct {
		Error       string `json:"error"`
		CleanedBody string `json:"cleaned_body"`
	}

	var ch chirp
	response := validateResponse{}

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

	if err := decoder.Decode(&ch); err != nil {
		response.Error = "Something went wrong"
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(ch.Body) < 140 {
		response.CleanedBody = purge_bad_words(ch.Body)
		w.WriteHeader(200)
	} else {
		response.Error = "Chirp is too long"
		w.WriteHeader(400)
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
