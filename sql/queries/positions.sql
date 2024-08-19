-- name: CreatePosition :one
insert into positions (id, created_at, updated_at, quantity, average_price, current_price, ppl, ticker)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *; 

-- name: GetLastPositionsTodayByTickerExcludingCurrent :many
select * from positions where ticker = $1 and id != $2 and date(created_at) = CURRENT_DATE order by created_at desc;

-- name: GetTodayPositions :many
select * from positions where date(created_at) = CURRENT_DATE order by created_at desc;

-- name: DeletePoistion :one
delete from positions where id = $1
returning *;

