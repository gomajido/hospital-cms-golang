CREATE TABLE IF NOT EXISTS articles (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    main_image VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    author_id CHAR(36) NOT NULL,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords TEXT,
    canonical_url VARCHAR(255),
    focus_keyphrase VARCHAR(255),
    og_title VARCHAR(255),
    og_description TEXT,
    og_image VARCHAR(255),
    UNIQUE KEY idx_articles_slug (slug),
    KEY idx_articles_status (status),
    KEY idx_articles_author (author_id),
    KEY idx_articles_published_at (published_at),
    CONSTRAINT fk_articles_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_articles_status CHECK (status IN ('published', 'draft', 'scheduled'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create a trigger to update updated_at timestamp
CREATE TRIGGER articles_updated_at
    BEFORE UPDATE ON articles
    FOR EACH ROW
    SET NEW.updated_at = CURRENT_TIMESTAMP;

-- Create fulltext search index for title and content
ALTER TABLE articles ADD FULLTEXT INDEX articles_fulltext_idx (title, content);
