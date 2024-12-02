-- Drop tables in correct order
DROP TABLE IF EXISTS article_categories;
DROP TABLE IF EXISTS categories;

-- Remove visitor_count column from articles
ALTER TABLE articles DROP COLUMN IF EXISTS visitor_count;
