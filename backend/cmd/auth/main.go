package main

import (
	"context"

	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/jackc/pgx/v5"

	"github.com/labstack/echo/v4"
)

func run(e *echo.Echo, ctx context.Context) error {
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

	auth.AddRoutes(e, config, queries)

	return e.Start(":8080")
}

func main() {
	e := echo.New()
	ctx := context.Background()

	if err := run(e, ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
