CREATE TABLE accesses (
    id UUID NOT NULL PRIMARY KEY,
    token VARCHAR(254) NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_unique_token ON accesses (token);

CREATE TRIGGER accesses_updated_at_modtime BEFORE UPDATE ON accesses FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();