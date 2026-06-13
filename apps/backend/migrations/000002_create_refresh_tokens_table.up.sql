CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,
    family_id UUID NOT NULL,
    user_id UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    replaced_by_token_id UUID,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_family_id ON refresh_tokens(family_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
