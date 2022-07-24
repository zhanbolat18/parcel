-- +goose Up
-- +goose StatementBegin
CREATE TABLE deliveries (
    id serial PRIMARY KEY,
    status VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    recipient_id VARCHAR(255) NOT NULL,
    courier_id VARCHAR(255) DEFAULT NULL,
    created_at VARCHAR(255) DEFAULT NULL,
    updated_at VARCHAR(255) DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE deliveries;
-- +goose StatementEnd
