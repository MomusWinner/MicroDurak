-- name: CreatePlayer :one
insert into player (name, age, rating)
values ($1, $2, 0)
returning id;

-- name: GetPlayerById :one
select * from player where id = $1;

-- name: UpdatePlayerRating :one
update player
   set rating = $2
 where id = $1
returning rating;

-- name: GetAllPlayers :many
select * from player;
