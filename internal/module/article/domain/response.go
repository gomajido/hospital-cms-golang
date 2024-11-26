package domain

// ArticleResponse represents the response for a single article
type ArticleResponse struct {
	Data *Article `json:"data"`
}

// ArticlesResponse represents the response for a list of articles
type ArticlesResponse struct {
	Data []Article     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PerPage     int   `json:"per_page"`
}
