-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS postgres_fdw;
CREATE TABLE public.logistics_orders_availability_shard_1 (
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
DROP TABLE public.logistics_orders_availability_shard_1;
-- +goose StatementEnd
