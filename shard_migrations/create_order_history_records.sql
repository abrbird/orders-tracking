-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS postgres_fdw;
CREATE TABLE order_history_records (
                                       id serial PRIMARY KEY,
                                       order_id bigint NOT NULL,
                                       status VARCHAR NOT NULL,
                                       confirmation VARCHAR NOT NULL,
                                       updated_at TIMESTAMP NOT NULL default current_timestamp,
                                       UNIQUE (order_id, status)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_history_records;
-- +goose StatementEnd