-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    uuid  UUID NOT NULL UNIQUE DEFAULT gen_random_uuid() PRIMARY KEY,
    amount           BIGINT NOT NULL DEFAULT 0 CHECK (amount >= 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
