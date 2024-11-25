CREATE TABLE IF NOT EXISTS user_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    token VARCHAR(255) NOT NULL,
    ability JSONB NOT NULL DEFAULT '[]',
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_user_tokens_token ON user_tokens(token) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_tokens_user_id ON user_tokens(user_id) WHERE deleted_at IS NULL;
