CREATE TABLE IF NOT EXISTS folders (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       TEXT NOT NULL,
    parent_id  UUID,
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_folders_id_owner UNIQUE (id, owner_id),

    CONSTRAINT fk_folders_parent
        FOREIGN KEY (parent_id, owner_id)
        REFERENCES folders (id, owner_id)
        ON DELETE CASCADE,

    CONSTRAINT uq_folder_name_in_parent UNIQUE NULLS NOT DISTINCT (name, parent_id, owner_id)
);
CREATE INDEX IF NOT EXISTS idx_folders_parent_id ON folders(parent_id);
CREATE INDEX IF NOT EXISTS idx_folders_owner_id  ON folders(owner_id);
