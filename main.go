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
	serveMux.HandleFunc("GET /healthz", healthz)

	serveMux.HandleFunc("GET /metrics", apiCfg.metrics)

	serveMux.HandleFunc("POST /reset", apiCfg.reset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
