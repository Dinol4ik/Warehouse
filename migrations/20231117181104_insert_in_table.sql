-- +goose Up
-- +goose StatementBegin
INSERT INTO public.item (code, name, size)
VALUES ('MP002XW0YADJ', 'Платье Sela', 'S'),
       ('MP002XW0L57W', 'Носки Calzedonia', 'S'),
       ('RTLACI608901', 'Платье Bad Queen', 'S');

INSERT INTO public.warehouse (uuid,name, is_available)
VALUES ('c8575d67-7cfb-48ff-b8ed-6f455a18cf05','EKB', true),
       ('3fec8d4d-a9fa-44f8-b8b7-20d79e74f00a','Tymen', true);

INSERT INTO public.item_warehouse (item_code, warehouse_id, amount, reserved)
VALUES ('MP002XW0YADJ', 'c8575d67-7cfb-48ff-b8ed-6f455a18cf05', 500, 200),
       ('MP002XW0L57W', '3fec8d4d-a9fa-44f8-b8b7-20d79e74f00a', 500, 200),
       ('RTLACI608901', 'c8575d67-7cfb-48ff-b8ed-6f455a18cf05', 500, 200);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
