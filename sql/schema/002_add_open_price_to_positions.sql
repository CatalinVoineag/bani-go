-- +goose Up
alter table positions
add column previous_close_price double precision;

-- +goose Down
alter table positions
drop column previous_close_price;
