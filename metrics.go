package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(fmt.Appendf(nil, "Hits: %v", apiCfg.fileserverHits.Load()))
}
