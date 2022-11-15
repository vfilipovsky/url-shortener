CREATE TABLE urls (
    id UUID NOT NULL PRIMARY KEY,
    access_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    url TEXT NOT NULL,
    pin VARCHAR(50) NOT NULL,
    is_secured BOOLEAN NOT NULL,
    alive_until TIMESTAMP NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_url_access FOREIGN KEY(access_id)
        REFERENCES accesses(id)
            ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_unique_code ON urls (code);

CREATE TRIGGER urls_updated_at_modtime BEFORE UPDATE ON urls FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();