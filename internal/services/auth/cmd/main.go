package main

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/config"
	"github.com/MommusWinner/MicroDurak/lib/validate"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/labstack/echo/v4"
)

func run(ctx context.Context, e *echo.Echo) error {
	config, err := config.Load()
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, config.DatabaseURL)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := database.New(conn)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	clientConn, err := grpc.NewClient(config.PlayersURL, opts...)
	if err != nil {
		return err
	}

	client := players.NewPlayersClient(clientConn)

	auth.AddRoutes(e, config, queries, client)

	return e.Start(":" + config.Port)
}

func main() {
	e := echo.New()
	ctx := context.Background()
	e.Validator = validate.NewHttpValidator(validator.New())

	if err := run(ctx, e); err != nil {
		e.Logger.Fatal(err)
	}
}
