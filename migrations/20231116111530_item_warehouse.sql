-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_warehouse
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_code    VARCHAR(12) NOT NULL REFERENCES "item" ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES "warehouse" on DELETE CASCADE,
    amount       INT  NOT NULL DEFAULT 0,
    reserved     INT NOT NULL DEFAULT 0
);

ALTER TABLE item_warehouse
    ADD CONSTRAINT item_code_warehouse_id UNIQUE (item_code, warehouse_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_warehouse;
-- +goose StatementEnd