package connection

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
)

type playerRepo struct {
	queries *database.Queries
}

func NewPlayerRepository(queries *database.Queries) *playerRepo {
	return &playerRepo{queries: queries}
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

	player := models.User{
		Id:     user.ID,
		Name:   user.Name,
		Age:    int(user.Age),
		Rating: int(user.Rating),
	}

	return &player, nil
}
