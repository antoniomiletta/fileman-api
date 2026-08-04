CREATE TABLE IF NOT EXISTS storage_cleanup_jobs (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    storage_key  TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    attempts     INT NOT NULL DEFAULT 0,
    last_error   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_storage_cleanup_jobs_pending
    ON storage_cleanup_jobs (created_at)
    WHERE status = 'pending';
