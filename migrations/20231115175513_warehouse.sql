-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS warehouse (
    uuid uuid primary key DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    is_available BOOL NOT NULL
);
CREATE UNIQUE INDEX warehouse_uuid_uindex
    ON warehouse (uuid);

CREATE UNIQUE INDEX warehouse_name_uindex
    ON warehouse (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS warehouse;
-- +goose StatementEnd
