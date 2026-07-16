CREATE TABLE IF NOT EXISTS token_blacklists (
    jti UUID PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_token_blacklists_expires_at ON token_blacklists (expires_at);
