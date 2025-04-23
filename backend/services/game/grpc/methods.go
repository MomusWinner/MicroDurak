package grpc

import (
	"context"

	"github.com/MommusWinner/MicroDurak/internal/game/v1"
	"github.com/MommusWinner/MicroDurak/services/game/config"
	"github.com/MommusWinner/MicroDurak/services/game/controller"
)

type GameServer struct {
	game.UnimplementedGameServer
	Config         *config.Config
	GameController *controller.GameController
}

func NewGameServer(gameController *controller.GameController, config *config.Config) *GameServer {
	return &GameServer{GameController: gameController, Config: config}
}

func (gs *GameServer) CreateGame(
	ctx context.Context,
	req *game.CreateGameRequest,
) (*game.CreateGameResponse, error) {
	createdGame, err := gs.GameController.CreateGame(req.UserIds)
	if err != nil {
		return nil, err
	}

	resp := &game.CreateGameResponse{GameId: createdGame.Id}
	return resp, err
}
