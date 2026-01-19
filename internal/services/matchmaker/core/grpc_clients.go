package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/contracts/game/v1"
	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/matchmaker/domain/infra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	playersClient     players.PlayersClient
	playersClientOnce sync.Once
	
	gameClient     game.GameClient
	gameClientOnce sync.Once
)

func makeSecureConnection(url string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakePlayersClient(cfg infra.Config) players.PlayersClient {
	playersClientOnce.Do(func() {
		conn, err := makeSecureConnection(cfg.GetPlayersURL())
		if err != nil {
			panic(fmt.Sprintf("could not create players client connection: %v", err))
		}
		playersClient = players.NewPlayersClient(conn)
		slog.Info("Players client initialized successfully.")
	})
	return playersClient
}

func MakeGameClient(cfg infra.Config) game.GameClient {
	gameClientOnce.Do(func() {
		conn, err := makeSecureConnection(cfg.GetGameURL())
		if err != nil {
			panic(fmt.Sprintf("could not create game client connection: %v", err))
		}
		gameClient = game.NewGameClient(conn)
		slog.Info("Game client initialized successfully.")
	})
	return gameClient
}
