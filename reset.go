package main

import (
	"context"
	"net/http"
)

func (apiCfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)
	apiCfg.db.ClearUsers(context.Background())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
