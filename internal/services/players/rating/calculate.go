package rating

import (
	"errors"
	"math"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
)

const kCoeff float64 = 30
const dCoeff float64 = 400

func CalculatePlayerScores(players []models.PlayerStats) ([]models.PlayerScore, error) {
	playerCount := len(players)
	if playerCount < 2 {
		return nil, errors.New("at least 2 players required")
	}

	playerScores := make([]models.PlayerScore, playerCount)
	var gamesCount float64 = float64(playerCount) * (float64(playerCount) - 1) / 2

	for i, playerStats := range players {
		var eTop float64 = 0
		for j, otherPlayerStats := range players {
			if i == j {
				continue
			}
			rDiff := float64(otherPlayerStats.Rating - playerStats.Rating)
			eTop += 1 / (1 + math.Pow(10, rDiff/dCoeff))
		}

		e := eTop / gamesCount
		s := float64(playerCount-playerStats.Place) / gamesCount

		ratingDelta := kCoeff * (s - e)
		newRating := playerStats.Rating + int(math.Floor(ratingDelta))
		playerScores[i] = models.PlayerScore{Id: playerStats.Id, Place: playerStats.Place, NewRating: newRating, RatingChange: int(math.Floor(ratingDelta))}
	}

	return playerScores, nil
}
