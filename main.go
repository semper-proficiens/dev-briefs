package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

const (
	customErrorMessage = "Oops. Something went wrong on our side."
)

type NewsItem struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Source      struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"source"`
	PublishedAt string `json:"publishedAt"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	http.HandleFunc("/", handler)
	http.HandleFunc("/news", newsHandler)
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	// if it's a htmx request
	if strings.Contains(r.Header.Get("HX-Request"), "true") {
		resp, err := http.Get("http://localhost:8081/api/everything-hacking-news")
		if err != nil {
			http.Error(w, customErrorMessage, http.StatusInternalServerError)
			return
		}
		defer func() {
			if err = resp.Body.Close(); err != nil {
				slog.Error("failed to close client connection", "newsHandler", err)
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read response body", http.StatusInternalServerError)
			return
		}

		var newsMap map[string]NewsItem
		err = json.Unmarshal(body, &newsMap)
		if err != nil {
			http.Error(w, "failed to unmarshal JSON", http.StatusInternalServerError)
			return
		}

		// Convert map to slice for easier template rendering
		var news []NewsItem
		for _, item := range newsMap {
			news = append(news, item)
		}

		tmpl := template.Must(template.ParseFiles("templates/news.html"))
		if err = tmpl.Execute(w, news); err != nil {
			slog.Error("failed to execute template", "newsHandler", err)
		}
	} else {
		// Redirect to the main page if accessed directly
		log.Println("Redirecting to home page", "newsHandler")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
