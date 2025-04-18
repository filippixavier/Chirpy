package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/filippixavier/Chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DB_URL")

	const filepathRoot = "."
	const port = "8080"

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
	}

	serveMux := http.NewServeMux()
	serveMux.Handle(
		"/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(
			http.FileServer(http.Dir(filepathRoot)),
		)),
	)
	serveMux.HandleFunc("GET /api/healthz", healthz)

	serveMux.HandleFunc("POST /api/validate_chirp", apiCfg.validate_chirp)

	serveMux.HandleFunc("POST /admin/reset", apiCfg.reset)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
