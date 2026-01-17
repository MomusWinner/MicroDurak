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
		id, err := uuid.Parse(placement.Id)
		if err != nil {
			uc.ctx.Logger().Error("Unprocessable player id", "player_id", placement.Id)
			return nil, ErrUnprocessableId
		}
		user, err := uc.ctx.Connection().UserRepository().GetById(ctx, id)
		if err != nil {
			uc.ctx.Logger().Error(err.Error())
			return nil, ErrInternal
		}
		if user == nil {
			uc.ctx.Logger().Error("Couldn't find player by id", "player_id", placement.Id)
			return nil, ErrPlayerNotFound
		}
		playerStats[i] = models.PlayerStats{
			Id:     id,
			Place:  placement.Place,
			Rating: user.Rating,
		}
	}

	playerRatings := make([]models.PlayerMatchResult, len(req.PlayerPlacements))

	err = uc.ctx.Connection().MatchRepository().WithTransaction(ctx,
		func(ctx context.Context, matchRepo repositories.MatchRepository, userRepo repositories.UserRepository) error {
			match, err := matchRepo.Add(ctx, len(req.PlayerPlacements), req.GameResult)
			if err != nil {
				return fmt.Errorf("Failed to create match: %w", err)
			}

			playerScores, err := rating.CalculatePlayerScores(playerStats)

			if err != nil {
				return fmt.Errorf("Failed to calculate ratings: %w", err)
			}

			for i, player := range playerScores {
				if err := matchRepo.AddPlayerToMatch(ctx, match.Id, player.Id, player.Place, int32(player.RatingChange)); err != nil {
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

func (uc *MatchUseCase) GetMatchResultById(ctx context.Context, req *props.GetMatchResultByIdReq) (resp *props.GetMatchResultByIdResp, err error) {
	match, err := uc.ctx.Connection().MatchRepository().GetById(ctx, req.Id)
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		return nil, ErrInternal
	}

	placements, err := uc.ctx.Connection().MatchRepository().GetPlayerPlacementsByMatchId(ctx, req.Id)
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		return nil, ErrInternal
	}

	resp = &props.GetMatchResultByIdResp{}
	resp.Match.Id = match.Id
	resp.Match.Result = models.GameResult_name[int32(match.GameResult)]

	for _, placement := range placements {
		playerDetail := models.PlayerMatchResultDetails{
			PlayerId:      placement.PlayerId,
			PlayerName:    placement.PlayerName,
			RatingChanged: placement.RatingChange,
			Place:         placement.PlayerPlace,
			CurrentRating: placement.CurrentRating,
		}
		resp.Match.Players = append(resp.Match.Players, playerDetail)
	}

	return resp, nil
}

func (uc *MatchUseCase) GetAllMatchResults(ctx context.Context, req *props.GetAllMatchResultsReq) (resp *props.GetAllMatchResultsResp, err error) {
	matches, err := uc.ctx.Connection().MatchRepository().GetAll(ctx)
	if err != nil {
		uc.ctx.Logger().Error(err.Error())
		return nil, ErrInternal
	}

	resp = &props.GetAllMatchResultsResp{}
	resp.Matches = make([]models.MatchDetails, len(matches))

	for i, match := range matches {
		resp.Matches[i].Id = match.Id
		resp.Matches[i].Result = models.GameResult_name[int32(match.GameResult)]

		placements, err := uc.ctx.Connection().MatchRepository().GetPlayerPlacementsByMatchId(ctx, match.Id)
		if err != nil {
			uc.ctx.Logger().Error(err.Error())
			continue
		}

		for _, placement := range placements {
			playerDetail := models.PlayerMatchResultDetails{
				PlayerId:      placement.PlayerId,
				PlayerName:    placement.PlayerName,
				RatingChanged: placement.RatingChange,
				Place:         placement.PlayerPlace,
				CurrentRating: placement.CurrentRating,
			}
			resp.Matches[i].Players = append(resp.Matches[i].Players, playerDetail)
		}
	}

	return resp, nil
}
