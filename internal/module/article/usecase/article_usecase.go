package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"

	"github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
)

type articleUsecase struct {
	articleRepo domain.ArticleRepository
}

// NewArticleUsecase creates a new instance of articleUsecase
func NewArticleUsecase(ar domain.ArticleRepository) domain.ArticleUsecase {
	return &articleUsecase{
		articleRepo: ar,
	}
}

func (u *articleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	article, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only increment visitor count for published articles
	if article.Status == "published" {
		// Use a goroutine to increment visitor count asynchronously
		go func() {
			// Create a new context with timeout for the async operation
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			if err := u.articleRepo.IncrementVisitorCount(ctx, id); err != nil {
				// Log error but don't affect the main response
				fmt.Printf("Error incrementing visitor count: %v\n", err)
			}
		}()
	}

	return article, nil
}

func (u *articleUsecase) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	article, err := u.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Only increment visitor count for published articles
	if article.Status == "published" {
		// Use a goroutine to increment visitor count asynchronously
		go func() {
			// Create a new context with timeout for the async operation
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			if err := u.articleRepo.IncrementVisitorCount(ctx, article.ID); err != nil {
				// Log error but don't affect the main response
				fmt.Printf("Error incrementing visitor count: %v\n", err)
			}
		}()
	}

	return article, nil
}

func (u *articleUsecase) List(ctx context.Context, page, limit int, status string) ([]domain.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Validate status if provided
	if status != "" {
		status = strings.ToLower(status)
		validStatuses := map[string]bool{
			"published": true,
			"draft":     true,
			"scheduled": true,
		}
		if !validStatuses[status] {
			return nil, 0, fmt.Errorf("invalid status: %s", status)
		}
	}

	return u.articleRepo.List(ctx, page, limit, status)
}

func (u *articleUsecase) Create(ctx context.Context, req domain.CreateArticleRequest) (*domain.Article, error) {
	article := &domain.Article{
		ID:              uuid.New(),
		Title:           req.Title,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		MainImage:       req.MainImage,
		Status:          req.Status,
		AuthorID:        req.AuthorID,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
		CanonicalURL:    req.CanonicalURL,
		FocusKeyphrase:  req.FocusKeyphrase,
		OGTitle:         req.OGTitle,
		OGDescription:   req.OGDescription,
		OGImage:         req.OGImage,
	}

	// Generate slug from title
	article.Slug = slug.Make(article.Title)

	// Validate status
	article.Status = strings.ToLower(article.Status)
	validStatuses := map[string]bool{
		"published": true,
		"draft":     true,
		"scheduled": true,
	}
	if !validStatuses[article.Status] {
		return nil, fmt.Errorf("invalid status: %s", article.Status)
	}

	// Set published_at for published articles
	if article.Status == "published" {
		now := time.Now()
		article.PublishedAt = &now
	}

	// Validate required fields
	if article.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if article.Content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if article.AuthorID == uuid.Nil {
		return nil, fmt.Errorf("author_id is required")
	}

	// Generate excerpt if not provided
	if article.Excerpt == "" {
		// Take first 160 characters of content (stripped of HTML)
		excerpt := stripHTML(article.Content)
		if len(excerpt) > 160 {
			excerpt = excerpt[:157] + "..."
		}
		article.Excerpt = excerpt
	}

	// Set SEO fields if not provided
	if article.MetaTitle == "" {
		article.MetaTitle = article.Title
	}
	if article.MetaDescription == "" {
		article.MetaDescription = article.Excerpt
	}
	if article.OGTitle == "" {
		article.OGTitle = article.Title
	}
	if article.OGDescription == "" {
		article.OGDescription = article.Excerpt
	}
	if article.OGImage == "" {
		article.OGImage = article.MainImage
	}

	if err := u.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

func (u *articleUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req domain.UpdateArticleRequest) (*domain.Article, error) {
	// Get existing article
	existing, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify user is the author
	if existing.AuthorID != userID {
		return nil, fmt.Errorf("unauthorized: only the author can update the article")
	}

	// Update fields if provided in request
	if req.Title != "" {
		existing.Title = req.Title
		existing.Slug = slug.Make(req.Title)
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.Excerpt != "" {
		existing.Excerpt = req.Excerpt
	}
	if req.MainImage != "" {
		existing.MainImage = req.MainImage
	}
	if req.Status != "" {
		status := strings.ToLower(req.Status)
		validStatuses := map[string]bool{
			"published": true,
			"draft":     true,
			"scheduled": true,
		}
		if !validStatuses[status] {
			return nil, fmt.Errorf("invalid status: %s", status)
		}
		existing.Status = status

		// Set published_at when status changes to published
		if status == "published" && existing.Status != "published" {
			now := time.Now()
			existing.PublishedAt = &now
		}
	}

	// Update SEO fields if provided
	if req.MetaTitle != "" {
		existing.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		existing.MetaDescription = req.MetaDescription
	}
	if req.MetaKeywords != "" {
		existing.MetaKeywords = req.MetaKeywords
	}
	if req.CanonicalURL != "" {
		existing.CanonicalURL = req.CanonicalURL
	}
	if req.FocusKeyphrase != "" {
		existing.FocusKeyphrase = req.FocusKeyphrase
	}
	if req.OGTitle != "" {
		existing.OGTitle = req.OGTitle
	}
	if req.OGDescription != "" {
		existing.OGDescription = req.OGDescription
	}
	if req.OGImage != "" {
		existing.OGImage = req.OGImage
	}

	// Generate excerpt if content changed and excerpt not provided
	if req.Content != "" && req.Excerpt == "" {
		excerpt := stripHTML(existing.Content)
		if len(excerpt) > 160 {
			excerpt = excerpt[:157] + "..."
		}
		existing.Excerpt = excerpt
	}

	if err := u.articleRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *articleUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Get existing article
	existing, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify user is the author
	if existing.AuthorID != userID {
		return fmt.Errorf("unauthorized: only the author can delete the article")
	}

	return u.articleRepo.Delete(ctx, id)
}

func (u *articleUsecase) IncrementVisitorCount(ctx context.Context, id uuid.UUID) error {
	// Check if article exists
	article, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if article == nil {
		return fmt.Errorf("article not found")
	}

	// Increment visitor count
	return u.articleRepo.IncrementVisitorCount(ctx, id)
}

// stripHTML removes HTML tags from a string
func stripHTML(html string) string {
	// This is a simple implementation. In production, you might want to use a proper HTML parser
	var sb strings.Builder
	var inTag bool
	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			sb.WriteRune(r)
		}
	}
	return strings.TrimSpace(sb.String())
}
