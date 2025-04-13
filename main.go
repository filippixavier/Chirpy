package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	serveMux.HandleFunc("/healthz", healthz)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
