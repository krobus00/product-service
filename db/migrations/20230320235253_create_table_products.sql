-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id varchar(36) UNIQUE,
    name text NOT NULL,
	description text NOT NULL,
    price float NOT NULL,
    thumbnail_id varchar(36) NOT NULL,
    owner_id varchar(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
