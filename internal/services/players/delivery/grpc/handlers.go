package grpc

import (
	"context"
	"errors"

	"github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/cases"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/props"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerService struct {
	players.UnimplementedPlayersServer
	ctx           domain.Context
	playerUseCase *cases.PlayerUseCase
	matchUseCase  *cases.MatchUseCase
}

func NewPlayerService(ctx domain.Context, playerUseCase *cases.PlayerUseCase, matchUseCase *cases.MatchUseCase) *PlayerService {
	return &PlayerService{ctx: ctx, playerUseCase: playerUseCase, matchUseCase: matchUseCase}
}

func (ps *PlayerService) CreatePlayer(ctx context.Context, req *players.CreatePlayerRequest) (*players.CreatePlayerReply, error) {
	resp, err := ps.playerUseCase.Create(props.CreatePlayerReq{Name: req.Name, Age: int(req.Age)})

	if err != nil {
		return nil, err
	}

	return &players.CreatePlayerReply{Id: resp.Id.String()}, nil
}

func (ps *PlayerService) GetPlayer(ctx context.Context, req *players.GetPlayerRequest) (*players.Player, error) {
	playerId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "player_id is not uuid").Err()
	}

	resp, err := ps.playerUseCase.GetById(props.GetPlayerByIdReq{Id: playerId})
	if err != nil {
		return nil, err
	}

	if resp.Player == nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}

	return &players.Player{
		Id:     resp.Player.Id.String(),
		Name:   resp.Player.Name,
		Age:    int32(resp.Player.Age),
		Rating: int32(resp.Player.Rating),
	}, nil
}

func grpcGameResultToDomain(gr players.GameResult) models.GameResult {
	switch gr {
	case players.GameResult_DRAW:
		return models.GameResult_DRAW
	case players.GameResult_WIN:
		return models.GameResult_WIN
	case players.GameResult_INTERRUPTED:
		return models.GameResult_INTERRUPTED
	default:
		panic("unknown game result")
	}
}

func (ps *PlayerService) CreateMatchResult(ctx context.Context, req *players.CreateMatchResultRequest) (*players.CreateMatchResultResponse, error) {
	gameResult := grpcGameResultToDomain(req.GameResult)

	playerPlacements := make([]models.PlayerPlacement, len(req.PlayerPlacements))
	for i, placement := range req.PlayerPlacements {
		playerPlacements[i] = models.PlayerPlacement{PlayerId: placement.PlayerId, PlayerPlace: int(placement.PlayerPlace)}
	}

	reqArgs := props.CreateMatchResutlReq{GameResult: gameResult, PlayerPlacements: playerPlacements}
	resp, err := ps.matchUseCase.CreateMatchResult(ctx, &reqArgs)
	if errors.Is(err, cases.ErrNoPlayers) {
		return nil, status.New(codes.InvalidArgument, "No players").Err()
	}

	playerRatings := make([]*players.PlayerPlacementResponse, len(resp.PlayerMatchResults))
	for i, rating := range resp.PlayerMatchResults {
		playerRatings[i] = &players.PlayerPlacementResponse{
			PlayerId:           rating.Id.String(),
			PlayerRating:       rating.Rating,
			PlayerRatingChange: rating.RatingChange,
		}
	}

	response := &players.CreateMatchResultResponse{
		MatchResultId: resp.MatchId.String(),
		PlayerRatings: playerRatings,
	}

	return response, nil
}
