-- +goose Up
alter table positions
add column previous_close_price bigint;

-- +goose Down
alter table positions
drop column previous_close_price;
