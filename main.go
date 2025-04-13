package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux := http.NewServeMux()
	serveMux.Handle(
		"/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(
			http.FileServer(http.Dir(filepathRoot)),
		)),
	)
	serveMux.HandleFunc("GET /api/healthz", healthz)

	serveMux.HandleFunc("POST /admin/reset", apiCfg.reset)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
