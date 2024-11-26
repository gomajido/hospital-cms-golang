package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Article represents the article entity
type Article struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Content         string     `json:"content"`
	Excerpt         string     `json:"excerpt"`
	MainImage       string     `json:"main_image"`
	Status          string     `json:"status"` // published, draft, scheduled
	AuthorID        uuid.UUID  `json:"author_id"`
	PublishedAt     *time.Time `json:"published_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	MetaTitle       string     `json:"meta_title"`
	MetaDescription string     `json:"meta_description"`
	MetaKeywords    string     `json:"meta_keywords"`
	CanonicalURL    string     `json:"canonical_url"`
	FocusKeyphrase  string     `json:"focus_keyphrase"`
	OGTitle         string     `json:"og_title"`
	OGDescription   string     `json:"og_description"`
	OGImage         string     `json:"og_image"`
}

// ArticleRepository defines the interface for article data operations
type ArticleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Article, error)
	GetBySlug(ctx context.Context, slug string) (*Article, error)
	List(ctx context.Context, page, limit int, status string) ([]Article, int64, error)
	Create(ctx context.Context, article *Article) error
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ArticleUsecase defines the interface for article business logic
type ArticleUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Article, error)
	GetBySlug(ctx context.Context, slug string) (*Article, error)
	List(ctx context.Context, page, limit int, status string) ([]Article, int64, error)
	Create(ctx context.Context, req CreateArticleRequest) (*Article, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req UpdateArticleRequest) (*Article, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
