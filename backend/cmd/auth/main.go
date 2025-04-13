package main

import (
	"context"

	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/jackc/pgx/v5"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	ctx := context.Background()

	config, err := config.Load()
	if err != nil {
		e.Logger.Fatal(err)
	}

	conn, err := pgx.Connect(ctx, config.DatabaseURL)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := database.New(conn)

	auth.AddRoutes(e, config, queries)

	e.Logger.Fatal(e.Start(":8080"))
}
