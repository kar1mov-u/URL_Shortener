-- +goose Up
CREATE TABLE urls(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_url TEXT NOT NULL,
    hashed_url TEXT NOT NULL, 
    created_at TIME NOT NULL,
    ttl TIME NOT NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE urls;
-- +goose StatementEnd
