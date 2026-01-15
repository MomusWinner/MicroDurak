package rating

import (
	"context"
	"errors"
	"math"

	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const kCoeff float64 = 30
const dCoeff float64 = 400

type PlayerScore struct {
	NewRating    int32
	RatingChange int
}

type PlayerStats struct {
	Rating int32
	Place  int32
}

type PlayerGetter interface {
	GetPlayerById(ctx context.Context, id pgtype.UUID) (database.Player, error)
}

func CalculatePlayerScores(
	ctx context.Context,
	dbQueries PlayerGetter,
	playerPlacements []*players.PlayerPlacementRequest,
) (map[string]PlayerScore, error) {
	playerCount := len(playerPlacements)
	if playerCount < 2 {
		return nil, errors.New("at least 2 players required")
	}

	playerScores := make(map[string]PlayerScore, playerCount)
	players := make(map[string]PlayerStats, playerCount)
	for _, playerPlacement := range playerPlacements {
		playerId, _ := uuid.Parse(playerPlacement.PlayerId)
		player, err := dbQueries.GetPlayerById(ctx, pgtype.UUID{Valid: true, Bytes: playerId})
		if err != nil {
			return nil, err
		}
		players[playerPlacement.PlayerId] = PlayerStats{player.Rating, playerPlacement.PlayerPlace}
	}

	var gamesCount float64 = float64(playerCount) * (float64(playerCount) - 1) / 2

	for playerId, playerStats := range players {
		var eTop float64 = 0
		for otherPlayerId, otherPlayerStats := range players {
			if playerId == otherPlayerId {
				continue
			}
			rDiff := float64(otherPlayerStats.Rating - playerStats.Rating)
			eTop += 1 / (1 + math.Pow(10, rDiff/dCoeff))
		}

		e := eTop / gamesCount
		s := float64(int32(playerCount)-playerStats.Place) / gamesCount

		ratingDelta := kCoeff * (s - e)
		newRating := playerStats.Rating + int32(math.Floor(ratingDelta))
		playerScores[playerId] = PlayerScore{newRating, int(math.Floor(ratingDelta))}
	}

	return playerScores, nil
}
