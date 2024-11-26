package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
)

type mysqlArticleRepository struct {
	db *sql.DB
}

func NewMySQLArticleRepository(db *sql.DB) domain.ArticleRepository {
	return &mysqlArticleRepository{
		db: db,
	}
}

// GetByID gets an article by ID
func (r *mysqlArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	article := &domain.Article{}
	query := `SELECT 
		id, title, slug, content, excerpt, main_image, status, author_id, 
		published_at, created_at, updated_at, meta_title, meta_description, 
		meta_keywords, canonical_url, focus_keyphrase, og_title, og_description, og_image 
		FROM articles WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID, &article.Title, &article.Slug, &article.Content,
		&article.Excerpt, &article.MainImage, &article.Status, &article.AuthorID,
		&article.PublishedAt, &article.CreatedAt, &article.UpdatedAt,
		&article.MetaTitle, &article.MetaDescription, &article.MetaKeywords,
		&article.CanonicalURL, &article.FocusKeyphrase, &article.OGTitle,
		&article.OGDescription, &article.OGImage,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found")
	}
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (r *mysqlArticleRepository) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	article := &domain.Article{}
	query := `SELECT 
		id, title, slug, content, excerpt, main_image, status, author_id, 
		published_at, created_at, updated_at, meta_title, meta_description, 
		meta_keywords, canonical_url, focus_keyphrase, og_title, og_description, og_image 
		FROM articles WHERE slug = ?`

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&article.ID, &article.Title, &article.Slug, &article.Content,
		&article.Excerpt, &article.MainImage, &article.Status, &article.AuthorID,
		&article.PublishedAt, &article.CreatedAt, &article.UpdatedAt,
		&article.MetaTitle, &article.MetaDescription, &article.MetaKeywords,
		&article.CanonicalURL, &article.FocusKeyphrase, &article.OGTitle,
		&article.OGDescription, &article.OGImage,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found")
	}
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (r *mysqlArticleRepository) List(ctx context.Context, page, limit int, status string) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64

	offset := (page - 1) * limit

	// Get total count
	countQuery := "SELECT COUNT(*) FROM articles"
	if status != "" {
		countQuery += " WHERE status = ?"
	}

	var err error
	if status != "" {
		err = r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	} else {
		err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	// Get articles
	query := `SELECT 
		id, title, slug, content, excerpt, main_image, status, author_id, 
		published_at, created_at, updated_at, meta_title, meta_description, 
		meta_keywords, canonical_url, focus_keyphrase, og_title, og_description, og_image 
		FROM articles`
	if status != "" {
		query += " WHERE status = ?"
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"

	var rows *sql.Rows
	if status != "" {
		rows, err = r.db.QueryContext(ctx, query, status, limit, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, query, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ID, &article.Title, &article.Slug, &article.Content,
			&article.Excerpt, &article.MainImage, &article.Status, &article.AuthorID,
			&article.PublishedAt, &article.CreatedAt, &article.UpdatedAt,
			&article.MetaTitle, &article.MetaDescription, &article.MetaKeywords,
			&article.CanonicalURL, &article.FocusKeyphrase, &article.OGTitle,
			&article.OGDescription, &article.OGImage,
		)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, article)
	}

	return articles, total, nil
}

func (r *mysqlArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	query := `INSERT INTO articles (
		id, title, slug, content, excerpt, main_image, status, author_id,
		published_at, created_at, updated_at, meta_title, meta_description,
		meta_keywords, canonical_url, focus_keyphrase, og_title, og_description, og_image
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		article.ID, article.Title, article.Slug, article.Content,
		article.Excerpt, article.MainImage, article.Status, article.AuthorID,
		article.PublishedAt, article.CreatedAt, article.UpdatedAt,
		article.MetaTitle, article.MetaDescription, article.MetaKeywords,
		article.CanonicalURL, article.FocusKeyphrase, article.OGTitle,
		article.OGDescription, article.OGImage,
	)

	return err
}

func (r *mysqlArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	query := `UPDATE articles SET
		title = ?, slug = ?, content = ?, excerpt = ?, main_image = ?,
		status = ?, author_id = ?, published_at = ?, updated_at = ?,
		meta_title = ?, meta_description = ?, meta_keywords = ?,
		canonical_url = ?, focus_keyphrase = ?, og_title = ?,
		og_description = ?, og_image = ?
		WHERE id = ?`

	article.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		article.Title, article.Slug, article.Content,
		article.Excerpt, article.MainImage, article.Status, article.AuthorID,
		article.PublishedAt, article.UpdatedAt,
		article.MetaTitle, article.MetaDescription, article.MetaKeywords,
		article.CanonicalURL, article.FocusKeyphrase, article.OGTitle,
		article.OGDescription, article.OGImage,
		article.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *mysqlArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM articles WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}
