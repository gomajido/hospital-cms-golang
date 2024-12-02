package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
)

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) domain.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

// GetByID gets an article by ID
func (r *articleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
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

// GetBySlug gets an article by slug
func (r *articleRepository) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
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

func (r *articleRepository) List(ctx context.Context, page, limit int, status string) ([]domain.Article, int64, error) {
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

	// Get articles with categories using LEFT JOIN
	query := `SELECT 
		a.id, a.title, a.slug, a.content, a.excerpt, a.main_image, a.status, a.author_id, 
		a.published_at, a.created_at, a.updated_at, a.meta_title, a.meta_description, 
		a.meta_keywords, a.canonical_url, a.focus_keyphrase, a.og_title, a.og_description, a.og_image,
		GROUP_CONCAT(
			DISTINCT JSON_OBJECT(
				'name', c.name,
				'slug', c.slug
			)
		) as categories_json
		FROM articles a
		LEFT JOIN article_categories ac ON a.id = ac.article_id
		LEFT JOIN categories c ON ac.category_id = c.id`
	if status != "" {
		query += " WHERE a.status = ?"
	}
	query += " GROUP BY a.id ORDER BY a.created_at DESC LIMIT ? OFFSET ?"

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
		var categoriesJSON sql.NullString
		err := rows.Scan(
			&article.ID, &article.Title, &article.Slug, &article.Content,
			&article.Excerpt, &article.MainImage, &article.Status, &article.AuthorID,
			&article.PublishedAt, &article.CreatedAt, &article.UpdatedAt,
			&article.MetaTitle, &article.MetaDescription, &article.MetaKeywords,
			&article.CanonicalURL, &article.FocusKeyphrase, &article.OGTitle,
			&article.OGDescription, &article.OGImage,
			&categoriesJSON,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse categories if they exist
		if categoriesJSON.Valid && categoriesJSON.String != "" {
			categoriesData := "[" + categoriesJSON.String + "]"
			err = json.Unmarshal([]byte(categoriesData), &article.Categories)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing categories: %v", err)
			}
		}

		articles = append(articles, article)
	}

	return articles, total, nil
}

func (r *articleRepository) Create(ctx context.Context, article *domain.Article) error {
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

func (r *articleRepository) Update(ctx context.Context, article *domain.Article) error {
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

func (r *articleRepository) Delete(ctx context.Context, id uuid.UUID) error {
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

func (r *articleRepository) IncrementVisitorCount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE articles 
		SET visitor_count = visitor_count + 1,
			updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

// GetCategories gets all categories for an article
func (r *articleRepository) GetCategories(ctx context.Context, articleID uuid.UUID) ([]domain.Category, error) {
	query := `SELECT c.id, c.name, c.slug, c.description, c.created_at, c.updated_at 
		FROM categories c 
		INNER JOIN article_categories ac ON c.id = ac.category_id 
		WHERE ac.article_id = ?`

	rows, err := r.db.QueryContext(ctx, query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategories updates the categories for an article
func (r *articleRepository) UpdateCategories(ctx context.Context, articleID uuid.UUID, categoryIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing article-category relationships
	_, err = tx.ExecContext(ctx, "DELETE FROM article_categories WHERE article_id = ?", articleID)
	if err != nil {
		return err
	}

	// Insert new article-category relationships
	for _, categoryID := range categoryIDs {
		_, err = tx.ExecContext(ctx, "INSERT INTO article_categories (article_id, category_id) VALUES (?, ?)",
			articleID, categoryID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
