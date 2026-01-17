-- name: CreateMatchResult :one
insert into match_result (player_count, game_result)
values ($1, $2)
returning *;

-- name: AddPlayerPlacement :one
insert into player_placement (match_result_id, player_id, player_place, rating_change)
values ($1, $2, $3, $4)
returning *;

-- name: GetMatchResultById :one
SELECT mr.id, mr.player_count, mr.game_result
FROM match_result mr
WHERE mr.id = $1;

-- name: GetPlayerPlacementsByMatchId :many
SELECT pp.player_id, pp.player_place, pp.rating_change, p.name, p.rating
FROM player_placement pp
JOIN player p ON pp.player_id = p.id
WHERE pp.match_result_id = $1
ORDER BY pp.player_place;

-- name: GetAllMatchResults :many
SELECT mr.id, mr.player_count, mr.game_result
FROM match_result mr
ORDER BY mr.id;
