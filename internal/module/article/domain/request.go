package domain

import (
	"time"

	"github.com/google/uuid"
)

// CreateArticleRequest represents the request to create a new article
type CreateArticleRequest struct {
	Title           string      `json:"title" validate:"required"`
	Content         string      `json:"content" validate:"required"`
	Excerpt         string      `json:"excerpt"`
	MainImage       string      `json:"main_image"`
	Status          string      `json:"status" validate:"required,oneof=published draft scheduled"`
	AuthorID        uuid.UUID   `json:"author_id" validate:"required"`
	CategoryIDs     []uuid.UUID `json:"category_ids" validate:"dive,required"`
	PublishedAt     *time.Time  `json:"published_at"`
	MetaTitle       string      `json:"meta_title"`
	MetaDescription string      `json:"meta_description"`
	MetaKeywords    string      `json:"meta_keywords"`
	CanonicalURL    string      `json:"canonical_url"`
	FocusKeyphrase  string      `json:"focus_keyphrase"`
	OGTitle         string      `json:"og_title"`
	OGDescription   string      `json:"og_description"`
	OGImage         string      `json:"og_image"`
}

// UpdateArticleRequest represents the request to update an existing article
type UpdateArticleRequest struct {
	Title           string      `json:"title"`
	Content         string      `json:"content"`
	Excerpt         string      `json:"excerpt"`
	MainImage       string      `json:"main_image"`
	Status          string      `json:"status" validate:"omitempty,oneof=published draft scheduled"`
	AuthorID        uuid.UUID   `json:"author_id"`
	CategoryIDs     []uuid.UUID `json:"category_ids" validate:"omitempty,dive,required"`
	PublishedAt     *time.Time  `json:"published_at"`
	MetaTitle       string      `json:"meta_title"`
	MetaDescription string      `json:"meta_description"`
	MetaKeywords    string      `json:"meta_keywords"`
	CanonicalURL    string      `json:"canonical_url"`
	FocusKeyphrase  string      `json:"focus_keyphrase"`
	OGTitle         string      `json:"og_title"`
	OGDescription   string      `json:"og_description"`
	OGImage         string      `json:"og_image"`
}

// ListArticlesRequest represents the request to list articles
type ListArticlesRequest struct {
	Page       int        `query:"page" validate:"omitempty,min=1"`
	Limit      int        `query:"limit" validate:"omitempty,min=1,max=100"`
	Status     string     `query:"status" validate:"omitempty,oneof=published draft scheduled"`
	CategoryID *uuid.UUID `query:"category_id"`
}
