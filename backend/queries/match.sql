-- name: CreateMatchResult :one
insert into match_result (player_count, game_result)
values ($1, $2)
returning *;

-- name: AddPlayerPlacement :one
insert into player_placement (match_result_id, player_id, player_place, rank_change)
values ($1, $2, $3, $4)
returning *;
