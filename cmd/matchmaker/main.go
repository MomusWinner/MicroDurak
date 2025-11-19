package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/game/v1"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/matchmaker"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/config"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/handlers"
	"github.com/MommusWinner/MicroDurak/services/matchmaker/types"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func run(ctx context.Context, e *echo.Echo) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config, err := config.Load()
	if err != nil {
		return err
	}

	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(opt)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	playerClientConn, err := grpc.NewClient(config.PlayersURL, opts...)
	if err != nil {
		return err
	}
	playerClient := players.NewPlayersClient(playerClientConn)

	gameClientConn, err := grpc.NewClient(config.GameURL, opts...)
	if err != nil {
		return err
	}
	gameClient := game.NewGameClient(gameClientConn)

	queueChan := make(chan types.MatchChan)
	cancelChan := make(chan types.MatchCancel)
	m := matchmaker.New(queueChan, cancelChan, config, redisClient, gameClient)

	handlers.AddRoutes(e, queueChan, cancelChan, config, playerClient)

	errChan := make(chan error, 2)

	go func() { errChan <- m.Start(ctx) }()
	go func() { errChan <- e.Start(":" + config.Port) }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	shutdownServices := func(shutdownErr error) error {
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			e.Logger.Errorf("Failed to shutdown Echo server: %v", err)
			// Return original error if exists, otherwise return shutdown error
			if shutdownErr != nil {
				return shutdownErr
			}
			return err
		}
		return shutdownErr
	}

	select {
	case err := <-errChan:
		return shutdownServices(err)
	case <-quit:
		e.Logger.Info("\nShutting down servers...")
		// Graceful shutdown without initial error
		return shutdownServices(nil)
	}
}

func main() {
	e := echo.New()
	ctx := context.Background()
	if err := run(ctx, e); err != nil {
		e.Logger.Fatal(err)
	}
}
