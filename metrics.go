package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetFileserverHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	log.Default().Println("fileserverHits reset")

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

	template, err := os.ReadFile("admin/index.html")
	if err != nil {
		log.Default().Printf("Failed to open file: %s", err)

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Error accessing the Admin Portal."))
		return
	}

	templateStr := string(template)
	updatedTemplate := fmt.Sprintf(templateStr, cfg.fileserverHits)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(updatedTemplate))
}
