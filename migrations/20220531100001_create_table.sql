-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE data
(
    id                     uuid        DEFAULT uuid_generate_v4() NOT NULL
        CONSTRAINT data_pk
            PRIMARY KEY,
    created_at             timestamptz DEFAULT NOW()              NOT NULL,
    account                text                                   NOT NULL,
    amount_cents           int8                                  NOT NULL,
    POS           text                                   NOT NULL,
    country         text                                   NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data;
-- +goose StatementEnd
