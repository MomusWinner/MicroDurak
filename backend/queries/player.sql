-- name: CreatePlayer :one
insert into player (name, age, rating)
values ($1, $2, 0)
returning id;

-- name: GetPlayerById :one
select * from player where id = $1;
