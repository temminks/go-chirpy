package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetFileserverHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = 0
		log.Default().Println("fileserverHits reset")

		w.Write([]byte("Hits reset to 0"))
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const port = "8080"
	var apiConfig apiConfig
	router := chi.NewRouter()

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	router.Get("/healthz", handlerReadiness)

	router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hits: %v", apiConfig.fileserverHits)))
	})

	router.Get("/reset", apiConfig.resetFileserverHits(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})).ServeHTTP)

	corsMux := middlewareCors(router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
