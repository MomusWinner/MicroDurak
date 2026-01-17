package cases

import (
	"context"
	"fmt"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/props"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/repositories"
	"github.com/MommusWinner/MicroDurak/internal/services/players/rating"
	"github.com/google/uuid"
)

type MatchUseCase struct {
	ctx domain.Context
}

func NewMatchUseCase(ctx domain.Context) *MatchUseCase {
	return &MatchUseCase{
		ctx: ctx,
	}
}

func (uc *MatchUseCase) CreateMatchResult(ctx context.Context, req *props.CreateMatchResutlReq) (resp *props.CreateMatchResutlResp, err error) {
	if len(req.PlayerPlacements) == 0 {
		err = ErrNoPlayers
		uc.ctx.Logger().Error(err.Error())
	}

	playerStats := make([]models.PlayerStats, len(req.PlayerPlacements))

	for i, placement := range req.PlayerPlacements {
		id, err := uuid.Parse(placement.PlayerId)
		if err != nil {
			uc.ctx.Logger().Error("Unprocessable player id", placement.PlayerId)
			return nil, ErrUnprocessableId
		}
		user, err := uc.ctx.Connection().UserRepository().GetById(ctx, id)
		if err != nil {
			uc.ctx.Logger().Error(err.Error())
			return nil, ErrInternal
		}
		if user == nil {
			uc.ctx.Logger().Error("Couldn't find player by id", placement.PlayerId)
			return nil, ErrPlayerNotFound
		}
		playerStats[i] = models.PlayerStats{
			Id:     id,
			Place:  placement.PlayerPlace,
			Rating: user.Rating,
		}
	}

	playerRatings := make([]models.PlayerMatchResult, len(req.PlayerPlacements))

	err = uc.ctx.Connection().MatchRepository().WithTransaction(ctx, func(ctx context.Context, matchRepo repositories.MatchRepository, userRepo repositories.UserRepository) error {
		match, err := matchRepo.Add(ctx, len(req.PlayerPlacements), req.GameResult)
		if err != nil {
			return fmt.Errorf("Failed to create match: %w", err)
		}

		playerScores, err := rating.CalculatePlayerScoresv2(playerStats)

		if err != nil {
			return fmt.Errorf("Failed to calculate ratings: %w", err)
		}

		for i, player := range playerScores {
			if err := matchRepo.AddPlayerToMatch(ctx, match.Id, player.Id, int32(player.RatingChange)); err != nil {
				return fmt.Errorf("Failed to add player to match: %w", err)
			}

			if err := userRepo.UpdatePlayerRating(ctx, player.Id, player.NewRating); err != nil {
				return fmt.Errorf("Failed to update player rating: %w", err)
			}

			playerRatings[i] = models.PlayerMatchResult{
				Id:           player.Id,
				Rating:       int32(player.NewRating),
				RatingChange: int32(player.RatingChange),
			}
		}

		resp = &props.CreateMatchResutlResp{
			MatchId:            match.Id,
			PlayerMatchResults: playerRatings,
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Match creation failed: %w", err)
	}

	return
}
