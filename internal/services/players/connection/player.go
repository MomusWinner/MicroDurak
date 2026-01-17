package connection

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type playerRepo struct {
	queries *database.Queries
	conn    *pgx.Conn
}

func NewPlayerRepository(queries *database.Queries, conn *pgx.Conn) *playerRepo {
	return &playerRepo{queries: queries, conn: conn}
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

func (r *playerRepo) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.conn.Query(ctx, "SELECT id, name, age, rating FROM player")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.User
	for rows.Next() {
		var p models.User
		var age int16
		var rating int32
		err := rows.Scan(&p.Id, &p.Name, &age, &rating)
		if err != nil {
			return nil, err
		}
		p.Age = int(age)
		p.Rating = int(rating)
		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}
