package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"
	var apiConfig apiConfig
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiConfig.resetFileserverHits)
	apiRouter.Post("/validate_chirp", handlerValidateChirp)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiConfig.handlerMetrics)

	// mount routers
	router.Mount("/api/", apiRouter)
	router.Mount("/admin", adminRouter)

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
