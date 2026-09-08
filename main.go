package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"golang.codes/handlers"
)

// The dev server. Production is a static build of the same pages — see
// cmd/build — so anything added here has to work as a plain file too.
func main() {
	// Unlocks the pages marked Private, which the static build leaves out.
	handlers.Dev = true

	if err := handlers.LoadTemplates("templates"); err != nil {
		log.Fatalf("Loading templates: %v", err)
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /steps", handlers.StepsHandler)
	mux.HandleFunc("GET /components", handlers.ComponentsHandler)
	mux.HandleFunc("GET /healthz", handlers.HealthHandler)

	// Everything else: an old URL that moved, or a 404.
	mux.HandleFunc("GET /", handlers.LegacyRedirectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Listening on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
