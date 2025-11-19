package rating

import (
	"context"
	"database/sql"
	"testing"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockDB struct {
	players map[string]*database.Player
}

func newMockDB() *mockDB {
	players := make(map[string]*database.Player)
	return &mockDB{players: players}
}

func (m *mockDB) addPlayer(playerId uuid.UUID, rating int32) {
	m.players[playerId.String()] = &database.Player{Rating: rating}
}

func (m *mockDB) GetPlayerById(ctx context.Context, id pgtype.UUID) (database.Player, error) {
	player := m.players[id.String()]
	if player == nil {
		return database.Player{}, sql.ErrNoRows
	}
	return *player, nil
}

func checkPlayerScores(t *testing.T, result map[string]PlayerScore, expected map[string]PlayerScore) {
	for id, exp := range expected {
		if result[id] != exp {
			t.Errorf("Player scores are not equal: \n"+
				"expected: %v\n"+
				"actual  : %v", exp, result[id])
		}
	}
}

func TestTwoPlayersSameRating(t *testing.T) {
	player1ID := uuid.New()
	player2ID := uuid.New()

	placements := []*players.PlayerPlacementRequest{
		{PlayerId: player1ID.String(), PlayerPlace: 1},
		{PlayerId: player2ID.String(), PlayerPlace: 2},
	}

	dbMock := newMockDB()
	dbMock.addPlayer(player1ID, 1000)
	dbMock.addPlayer(player2ID, 1000)

	result, err := CalculatePlayerScores(context.Background(), dbMock, placements)
	if err != nil {
		t.Error("Returned error, but should not")
	}

	expected := make(map[string]PlayerScore)
	expected[player1ID.String()] = PlayerScore{NewRating: 1015, RatingChange: 15}
	expected[player2ID.String()] = PlayerScore{NewRating: 985, RatingChange: -15}
	checkPlayerScores(t, result, expected)
}

func TestOnePlayer(t *testing.T) {
	player1ID := uuid.New()

	placements := []*players.PlayerPlacementRequest{
		{PlayerId: player1ID.String(), PlayerPlace: 1},
	}

	dbMock := newMockDB()
	dbMock.addPlayer(player1ID, 1000)

	result, err := CalculatePlayerScores(context.Background(), dbMock, placements)
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

	placements := []*players.PlayerPlacementRequest{
		{PlayerId: player1ID.String(), PlayerPlace: 1},
		{PlayerId: player2ID.String(), PlayerPlace: 2},
		{PlayerId: player3ID.String(), PlayerPlace: 3},
	}

	dbMock := newMockDB()
	dbMock.addPlayer(player1ID, 1000)
	dbMock.addPlayer(player2ID, 1000)
	dbMock.addPlayer(player3ID, 1000)

	result, err := CalculatePlayerScores(context.Background(), dbMock, placements)
	if err != nil {
		t.Error("Returned error, but should not")
	}

	expected := make(map[string]PlayerScore)
	expected[player1ID.String()] = PlayerScore{NewRating: 1010, RatingChange: 10}
	expected[player2ID.String()] = PlayerScore{NewRating: 1000, RatingChange: 0}
	expected[player3ID.String()] = PlayerScore{NewRating: 990, RatingChange: -10}
	checkPlayerScores(t, result, expected)
}

func TestDraw(t *testing.T) {
	player1ID := uuid.New()
	player2ID := uuid.New()
	player3ID := uuid.New()

	placements := []*players.PlayerPlacementRequest{
		{PlayerId: player1ID.String(), PlayerPlace: 1},
		{PlayerId: player2ID.String(), PlayerPlace: 2},
		{PlayerId: player3ID.String(), PlayerPlace: 2},
	}

	dbMock := newMockDB()
	dbMock.addPlayer(player1ID, 1000)
	dbMock.addPlayer(player2ID, 1000)
	dbMock.addPlayer(player3ID, 1000)

	result, err := CalculatePlayerScores(context.Background(), dbMock, placements)
	if err != nil {
		t.Error("Returned error, but should not")
	}

	expected := make(map[string]PlayerScore)
	expected[player1ID.String()] = PlayerScore{NewRating: 1010, RatingChange: 10}
	expected[player2ID.String()] = PlayerScore{NewRating: 1000, RatingChange: 0}
	expected[player3ID.String()] = PlayerScore{NewRating: 1000, RatingChange: 0}
	checkPlayerScores(t, result, expected)
}
