CREATE TABLE IF NOT EXISTS user_tokens (
    id CHAR(36) NOT NULL COMMENT 'UUID v4',
    user_id CHAR(36) NOT NULL,
    token VARCHAR(255) NOT NULL,
    ability JSON NOT NULL DEFAULT ('[]'),
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_tokens_token (token),
    KEY idx_user_tokens_user_id (user_id),
    KEY idx_user_tokens_expired_at (expired_at),
    KEY idx_user_tokens_deleted_at (deleted_at),
    CONSTRAINT fk_user_tokens_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
