-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY,
    name TEXT,
    quantity INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id)
);

CREATE TRIGGER IF NOT EXISTS trigger_items_updated_at
    AFTER UPDATE ON items
    FOR EACH ROW
BEGIN
    UPDATE items SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;

DROP TRIGGER IF EXISTS trigger_items_updated_at;
-- +goose StatementEnd
