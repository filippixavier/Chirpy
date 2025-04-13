package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func (apiCfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	raw, err := os.ReadFile("./metrics.html")

	if err != nil {
		fmt.Println("Error when trying to read file")
		w.Write(fmt.Appendf(nil, "Hits: %v", apiCfg.fileserverHits.Load()))
		return
	}

	template := string(raw)

	template = strings.Replace(template, "%d", fmt.Sprint(apiCfg.fileserverHits.Load()), 1)

	w.Write([]byte(template))
}
