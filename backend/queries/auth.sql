-- name: CreateAuth :one
insert into player_auth (player_id, email, password)
values ($1, $2, $3)
returning id;

-- name: GetAuthByEmail :one
select * from player_auth
 where email = $1
 limit 1;

-- name: CheckEmail :one
select count(*) from player_auth
 where email = $1
 limit 1;
