-- Drop the fulltext index
ALTER TABLE articles DROP INDEX articles_fulltext_idx;

-- Drop the trigger
DROP TRIGGER IF EXISTS articles_updated_at;

-- Drop the table
DROP TABLE IF EXISTS articles;
