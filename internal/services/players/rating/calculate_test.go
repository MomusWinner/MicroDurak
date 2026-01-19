package rating

import (
	"testing"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
)

func checkPlayerScores(t *testing.T, result []models.PlayerScore, expected map[uuid.UUID]models.PlayerScore) {
	if len(result) != len(expected) {
		t.Errorf("Result length mismatch: expected %d, got %d", len(expected), len(result))
		return
	}

	for _, score := range result {
		exp, ok := expected[score.Id]
		if !ok {
			t.Errorf("Unexpected player in result: %v", score.Id)
			continue
		}

		if score.NewRating != exp.NewRating {
			t.Errorf("Player %v: NewRating mismatch - expected %d, got %d", score.Id, exp.NewRating, score.NewRating)
		}
		if score.RatingChange != exp.RatingChange {
			t.Errorf("Player %v: RatingChange mismatch - expected %d, got %d", score.Id, exp.RatingChange, score.RatingChange)
		}
		if score.Place != exp.Place {
			t.Errorf("Player %v: Place mismatch - expected %d, got %d", score.Id, exp.Place, score.Place)
		}
	}
}

func TestTwoPlayersSameRating(t *testing.T) {
	player1ID := uuid.New()
	player2ID := uuid.New()

	players := []models.PlayerStats{
		{Id: player1ID, Rating: 1000, Place: 1},
		{Id: player2ID, Rating: 1000, Place: 2},
	}

	result, err := CalculatePlayerScores(players)
	if err != nil {
		t.Error("Returned error, but should not:", err)
		return
	}

	expected := make(map[uuid.UUID]models.PlayerScore)
	expected[player1ID] = models.PlayerScore{Id: player1ID, Place: 1, NewRating: 1015, RatingChange: 15}
	expected[player2ID] = models.PlayerScore{Id: player2ID, Place: 2, NewRating: 985, RatingChange: -15}
	checkPlayerScores(t, result, expected)
}

func TestOnePlayer(t *testing.T) {
	player1ID := uuid.New()

	players := []models.PlayerStats{
		{Id: player1ID, Rating: 1000, Place: 1},
	}

	result, err := CalculatePlayerScores(players)
	if err == nil {
		t.Error("No error, but should be")
	}

	if result != nil {
		t.Error("Returned result, should be nil")
	}
}

func TestManyPlayers(t *testing.T) {
	player1ID := uuid.New()
	player2ID := uuid.New()
	player3ID := uuid.New()

	players := []models.PlayerStats{
		{Id: player1ID, Rating: 1000, Place: 1},
		{Id: player2ID, Rating: 1000, Place: 2},
		{Id: player3ID, Rating: 1000, Place: 3},
	}

	result, err := CalculatePlayerScores(players)
	if err != nil {
		t.Error("Returned error, but should not:", err)
		return
	}

	expected := make(map[uuid.UUID]models.PlayerScore)
	expected[player1ID] = models.PlayerScore{Id: player1ID, Place: 1, NewRating: 1010, RatingChange: 10}
	expected[player2ID] = models.PlayerScore{Id: player2ID, Place: 2, NewRating: 1000, RatingChange: 0}
	expected[player3ID] = models.PlayerScore{Id: player3ID, Place: 3, NewRating: 990, RatingChange: -10}
	checkPlayerScores(t, result, expected)
}

func TestDraw(t *testing.T) {
	player1ID := uuid.New()
	player2ID := uuid.New()
	player3ID := uuid.New()

	players := []models.PlayerStats{
		{Id: player1ID, Rating: 1000, Place: 1},
		{Id: player2ID, Rating: 1000, Place: 2},
		{Id: player3ID, Rating: 1000, Place: 2},
	}

	result, err := CalculatePlayerScores(players)
	if err != nil {
		t.Error("Returned error, but should not:", err)
		return
	}

	expected := make(map[uuid.UUID]models.PlayerScore)
	expected[player1ID] = models.PlayerScore{Id: player1ID, Place: 1, NewRating: 1010, RatingChange: 10}
	expected[player2ID] = models.PlayerScore{Id: player2ID, Place: 2, NewRating: 1000, RatingChange: 0}
	expected[player3ID] = models.PlayerScore{Id: player3ID, Place: 2, NewRating: 1000, RatingChange: 0}
	checkPlayerScores(t, result, expected)
}
