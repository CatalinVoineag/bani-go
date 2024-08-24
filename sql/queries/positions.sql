-- name: CreatePosition :one
insert into positions (id, created_at, updated_at, quantity, average_price, current_price, ppl, ticker)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *; 

-- name: GetLastPositionTodayByTickerExcludingCurrent :one
select * from positions where ticker = $1 and id != $2 and date(created_at) = CURRENT_DATE order by created_at desc limit 1;

-- name: GetLastPositionTodayByTicker :one
select * from positions where ticker = $1 and date(created_at) = CURRENT_DATE order by created_at desc limit 1;

-- name: GetTodayPositions :many
select * from positions where date(created_at) = CURRENT_DATE order by created_at desc;

-- name: GetTodayPositionsTickers :many
select distinct(ticker, id) from positions where date(created_at) = CURRENT_DATE;

-- name: UpdatePosition :one
update positions set quantity = $1, average_price = $2, current_price = $3, ppl = $4, ticker = $5 where id = $6
returning *; 

-- name: UpdatePreviousClosedPrice :one
update positions set previous_close_price = $1 where id = $2
returning *;

-- name: DeletePoistion :one
delete from positions where id = $1
returning *;

