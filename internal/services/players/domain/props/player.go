package props

import (
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/google/uuid"
)

type CreatePlayerReq struct {
	Name string
	Age  int
}

type CreatePlayerResp struct {
	Id uuid.UUID
}

type GetPlayerByIdReq struct {
	Id uuid.UUID
}

type GetPlayerByIdResp struct {
	Player *models.User
}

type GetAllPlayersReq struct {
}

type GetAllPlayersResp struct {
	Players []models.User
}

type CreateMatchResutlReq struct {
	GameResult       models.GameResult
	PlayerPlacements []models.PlayerPlacement
}

type CreateMatchResutlResp struct {
	MatchId            uuid.UUID
	PlayerMatchResults []models.PlayerMatchResult
}

type GetMatchResultByIdReq struct {
	Id uuid.UUID
}

type GetMatchResultByIdResp struct {
	Match models.MatchDetails
}

type GetAllMatchResultsReq struct{}

type GetAllMatchResultsResp struct {
	Matches []models.MatchDetails
}
