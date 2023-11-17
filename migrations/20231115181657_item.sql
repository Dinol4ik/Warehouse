-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item
(
    code   VARCHAR(12) PRIMARY KEY,
    name   TEXT NOT NULL,
    size   TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists item
-- +goose StatementEnd
