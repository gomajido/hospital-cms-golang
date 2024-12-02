-- Add register_from and is_active columns to users table
ALTER TABLE users
ADD COLUMN register_from VARCHAR(50) NOT NULL DEFAULT 'web',
ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 0;
