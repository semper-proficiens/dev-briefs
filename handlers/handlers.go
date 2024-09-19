package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/semper-proficiens/dev-briefs/components"
	"github.com/semper-proficiens/dev-briefs/types"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

const (
	customErrorMessage = "Oops. Something went wrong on our side."
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {
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

		var newsMap map[string]types.NewsItem
		err = json.Unmarshal(body, &newsMap)
		if err != nil {
			http.Error(w, "failed to unmarshal JSON", http.StatusInternalServerError)
			return
		}

		// Convert map to slice for easier template rendering
		var news []types.NewsItem
		for _, item := range newsMap {
			news = append(news, item)
		}

		if err = components.News(news).Render(r.Context(), w); err != nil {
			http.Error(w, customErrorMessage, http.StatusInternalServerError)
			return
		}
	} else {
		// Redirect to the main page if accessed directly
		log.Println("Redirecting to home page", "newsHandler")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func CollapseDivHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("HX-Request"), "true") {
		id := chi.URLParam(r, "id")
		log.Println("Grabbed collapse id:", id)
		components.EmptyDiv(id).Render(r.Context(), w)
	} else {
		// Redirect to the main page if accessed directly
		log.Println("Redirecting to home page", "newsHandler")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
