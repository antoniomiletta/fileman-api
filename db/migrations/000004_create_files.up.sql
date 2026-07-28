CREATE TABLE IF NOT EXISTS files (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         TEXT NOT NULL,
    parent_id    UUID NOT NULL,
    owner_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    size         BIGINT NOT NULL DEFAULT 0,
    mime_type    TEXT NOT NULL DEFAULT 'application/octet-stream',
    storage_key  TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_files_parent
        FOREIGN KEY (parent_id, owner_id)
        REFERENCES folders (id, owner_id)
        ON DELETE CASCADE,

    CONSTRAINT uq_file_name_in_parent UNIQUE (name, parent_id)
);

CREATE INDEX IF NOT EXISTS idx_files_parent_id ON files(parent_id);
CREATE INDEX IF NOT EXISTS idx_files_owner_id  ON files(owner_id);
