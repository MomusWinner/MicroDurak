package connection

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/models"
	"github.com/google/uuid"
)

type authRepo struct {
	queries *database.Queries
}

func NewAuthRepository(queries *database.Queries) *authRepo {
	return &authRepo{queries: queries}
}

func (r *authRepo) Add(auth *models.AuthUser) error {
	_, err := r.queries.CreateAuth(context.TODO(), database.CreateAuthParams{PlayerID: auth.PlayerId, Email: auth.Email, Password: auth.Password})
	return err
}

func (r *authRepo) GetByEmail(email string) (*models.AuthUser, error) {
	auth, err := r.queries.GetAuthByEmail(context.TODO(), email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	domainAuth := databaseAuthToDomain(auth)
	return &domainAuth, nil
}

func (r *authRepo) Delete(id uuid.UUID) error {
	panic("Not implemented")
}

func databaseAuthToDomain(dbPlayer database.PlayerAuth) models.AuthUser {
	return models.AuthUser{
		Id:       dbPlayer.ID,
		PlayerId: dbPlayer.PlayerID,
		Email:    dbPlayer.Email,
		Password: dbPlayer.Password,
	}
}
