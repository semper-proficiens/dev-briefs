package types

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
