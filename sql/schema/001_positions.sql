-- +goose Up
create table positions(
  id UUID primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  quantity double precision not null,
  average_price double precision not null,
  current_price double precision not null,
  ppl double precision not null,
  ticker text not null
);
-- +goose Down
drop table positions;
