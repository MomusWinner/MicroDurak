package connection

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type playerRepo struct {
	queries *database.Queries
	pool    *pgxpool.Pool
}

func NewPlayerRepository(queries *database.Queries, pool *pgxpool.Pool) *playerRepo {
	return &playerRepo{queries: queries, pool: pool}
}

func (r *playerRepo) Add(ctx context.Context, player *models.User) (uuid.UUID, error) {
	id, err := r.queries.CreatePlayer(ctx, database.CreatePlayerParams{Name: player.Name, Age: int16(player.Age)})
	return id, err
}

func (r *playerRepo) UpdatePlayerRating(ctx context.Context, playerID uuid.UUID, newRating int) error {
	_, err := r.queries.UpdatePlayerRating(ctx, database.UpdatePlayerRatingParams{ID: playerID, Rating: int32(newRating)})
	return err
}

func (r *playerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}

func (r *playerRepo) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := r.queries.GetPlayerById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	player := databasePlayerToDomain(user)
	return &player, nil
}

func (r *playerRepo) GetAll(ctx context.Context) ([]models.User, error) {
	players, err := r.queries.GetAllPlayers(ctx)
	if err != nil {
		return nil, err
	}

	return databasePlayersToDomain(players), nil
}

func databasePlayerToDomain(player database.Player) models.User {
	return models.User{
		Id:     player.ID,
		Name:   player.Name,
		Age:    int(player.Age),
		Rating: int(player.Rating),
	}
}

func databasePlayersToDomain(players []database.Player) []models.User {
	users := make([]models.User, len(players))

	for i, player := range players {
		users[i] = databasePlayerToDomain(player)
	}

	return users
}
