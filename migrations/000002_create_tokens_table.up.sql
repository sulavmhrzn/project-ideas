CREATE TABLE IF NOT EXISTS tokens(
    token text PRIMARY KEY,
    userId int NOT NULL REFERENCES users ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    scope text NOT NULL
);