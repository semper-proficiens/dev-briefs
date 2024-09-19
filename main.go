package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/semper-proficiens/dev-briefs/handlers"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := chi.NewRouter()

	// Serve static files from the "assets" directory
	fs := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	r.Get("/", handlers.HomeHandler)
	r.Get("/news", handlers.NewsHandler)
	r.Get("/collapse/{id}", handlers.CollapseDivHandler)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
