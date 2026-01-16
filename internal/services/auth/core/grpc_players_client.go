package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/infra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client     players.PlayersClient
	clientOnce sync.Once
)

func makeSecureConnection(cfg infra.Config) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.GetPlayersURL(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not connect to grpc client: %v", err)
	}
	return conn, nil
}

func MakePlayersClient(cfg infra.Config) players.PlayersClient {
	clientOnce.Do(func() {
		conn, err := makeSecureConnection(cfg)
		if err != nil {
			panic(fmt.Sprintf("could not create client connection: %v", err))
		}
		client = players.NewPlayersClient(conn)
		slog.Info("Parking client initialized successfully.")
	})
	return client
}

