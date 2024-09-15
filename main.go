package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
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
	http.HandleFunc("/", handler)
	http.HandleFunc("/content", contentHandler)
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8081/api/everything-hacking-news")
	if err != nil {
		http.Error(w, "Failed to fetch news", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var newsMap map[string]NewsItem
	err = json.Unmarshal(body, &newsMap)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Convert map to slice for easier template rendering
	var news []NewsItem
	for _, item := range newsMap {
		news = append(news, item)
	}

	tmpl := template.Must(template.ParseFiles("templates/content.html"))
	tmpl.Execute(w, news)
}
