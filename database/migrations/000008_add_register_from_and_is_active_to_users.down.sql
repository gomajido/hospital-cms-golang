-- Remove register_from and is_active columns from users table
ALTER TABLE users
DROP COLUMN register_from,
DROP COLUMN is_active;
