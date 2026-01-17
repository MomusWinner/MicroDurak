package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/cases"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/models"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/props"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PlayerHandler struct {
	ctx           domain.Context
	playerUseCase *cases.PlayerUseCase
	matchUseCase  *cases.MatchUseCase
}

func NewPlayerHandler(ctx domain.Context, playerUseCase *cases.PlayerUseCase, matchUseCase *cases.MatchUseCase) *PlayerHandler {
	return &PlayerHandler{
		ctx:           ctx,
		playerUseCase: playerUseCase,
		matchUseCase:  matchUseCase,
	}
}

type PlayerResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Rating int    `json:"rating"`
}

type PlayersResponse struct {
	Players []PlayerResponse `json:"players"`
}

type CreateMatchRequest struct {
	GameResult       string               `json:"game_result"`
	PlayerPlacements []PlayerPlacementReq `json:"player_placements"`
}

type PlayerPlacementReq struct {
	PlayerId    string `json:"player_id" validate:"required,uuid"`
	PlayerPlace int    `json:"player_place" validate:"required,min=1"`
}

type PlayerMatchResultResponse struct {
	Id           string `json:"id"`
	Rating       int32  `json:"rating"`
	RatingChange int32  `json:"rating_change"`
}

type CreateMatchResponse struct {
	MatchId            string                      `json:"match_id"`
	PlayerMatchResults []PlayerMatchResultResponse `json:"player_match_results"`
}

var internalServerError = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")

// GetAll retrieves a list of all players
// @Summary Get all players
// @Description Returns a list of all players in the system
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {object} PlayersResponse "List of players"
// @Failure 500 "Internal server error"
// @Router /players [get]
func (h *PlayerHandler) GetAll(c echo.Context) error {
	resp, err := h.playerUseCase.GetAll(props.GetAllPlayersReq{})

	if err != nil {
		if errors.Is(err, cases.ErrNoPlayers) {
			return c.JSON(http.StatusOK, PlayersResponse{Players: []PlayerResponse{}})
		}
		return internalServerError
	}

	players := make([]PlayerResponse, 0, len(resp.Players))
	for _, p := range resp.Players {
		players = append(players, PlayerResponse{
			Id:     p.Id.String(),
			Name:   p.Name,
			Age:    p.Age,
			Rating: p.Rating,
		})
	}

	return c.JSON(http.StatusOK, PlayersResponse{Players: players})
}

// GetById retrieves a player by ID
// @Summary Get player by ID
// @Description Returns player information by their unique identifier
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player UUID" format(uuid)
// @Success 200 {object} PlayerResponse "Player information"
// @Failure 400 "Invalid ID format"
// @Failure 404 "Player not found"
// @Failure 500 "Internal server error"
// @Router /players/{id} [get]
func (h *PlayerHandler) GetById(c echo.Context) error {
	idParam := c.Param("id")
	playerId, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid player ID format")
	}

	resp, err := h.playerUseCase.GetById(props.GetPlayerByIdReq{Id: playerId})

	if err != nil {
		return internalServerError
	}

	if resp.Player == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Player not found")
	}

	return c.JSON(http.StatusOK, PlayerResponse{
		Id:     resp.Player.Id.String(),
		Name:   resp.Player.Name,
		Age:    resp.Player.Age,
		Rating: resp.Player.Rating,
	})
}

// CreateMatch creates a match result and updates player ratings
// @Summary Create match result
// @Description Creates a match result with player placements and updates player ratings based on the game outcome
// @Tags matches
// @Accept json
// @Produce json
// @Param request body CreateMatchRequest true "Match result data"
// @Success 201 {object} CreateMatchResponse "Match result created successfully"
// @Failure 400 "Bad request - validation error"
// @Failure 500 "Internal server error"
// @Router /matches [post]
func (h *PlayerHandler) CreateMatch(c echo.Context) error {
	req := new(CreateMatchRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := c.Validate(req); err != nil {
		h.ctx.Logger().Info("Hello")
		h.ctx.Logger().Info(err.Error())
		return err
	}

	var gameResult models.GameResult
	switch req.GameResult {
	case "win":
		gameResult = models.GameResult_WIN
	case "draw":
		gameResult = models.GameResult_DRAW
	case "interrupted":
		gameResult = models.GameResult_INTERRUPTED
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game_result value. Must be win, draw, or interrupted")
	}

	playerPlacements := make([]models.PlayerPlacement, len(req.PlayerPlacements))
	for i, placement := range req.PlayerPlacements {
		playerPlacements[i] = models.PlayerPlacement{
			Id:    placement.PlayerId,
			Place: placement.PlayerPlace,
		}
	}

	matchReq := &props.CreateMatchResutlReq{
		GameResult:       gameResult,
		PlayerPlacements: playerPlacements,
	}

	resp, err := h.matchUseCase.CreateMatchResult(context.Background(), matchReq)
	if err != nil {
		if errors.Is(err, cases.ErrNoPlayers) {
			return echo.NewHTTPError(http.StatusBadRequest, "No players provided")
		}
		if errors.Is(err, cases.ErrUnprocessableId) {
			return echo.NewHTTPError(http.StatusBadRequest, "Unprocessable player id")
		}
		if errors.Is(err, cases.ErrPlayerNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Couldn't find the player")
		}
		return internalServerError
	}

	playerResults := make([]PlayerMatchResultResponse, len(resp.PlayerMatchResults))
	for i, result := range resp.PlayerMatchResults {
		playerResults[i] = PlayerMatchResultResponse{
			Id:           result.Id.String(),
			Rating:       result.Rating,
			RatingChange: result.RatingChange,
		}
	}

	return c.JSON(http.StatusCreated, CreateMatchResponse{
		MatchId:            resp.MatchId.String(),
		PlayerMatchResults: playerResults,
	})
}

// GetMatchResultById retrieves a match result by ID
// @Summary Get match result by ID
// @Description Returns match result information by its unique identifier
// @Tags matches
// @Accept json
// @Produce json
// @Param id path string true "Match UUID" format(uuid)
// @Success 200 {object} models.MatchDetails "Match result information"
// @Failure 400 "Invalid ID format"
// @Failure 404 "Match not found"
// @Failure 500 "Internal server error"
// @Router /matches/{id} [get]
func (h *PlayerHandler) GetMatchResultById(c echo.Context) error {
	idParam := c.Param("id")
	matchId, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid match ID format")
	}

	resp, err := h.matchUseCase.GetMatchResultById(context.Background(), &props.GetMatchResultByIdReq{Id: matchId})
	if err != nil {
		return internalServerError
	}

	return c.JSON(http.StatusOK, resp.Match)
}

// GetAllMatchResults retrieves all match results
// @Summary Get all match results
// @Description Returns a list of all match results in the system
// @Tags matches
// @Accept json
// @Produce json
// @Success 200 {object} props.GetAllMatchResultsResp "List of match results"
// @Failure 500 "Internal server error"
// @Router /matches [get]
func (h *PlayerHandler) GetAllMatchResults(c echo.Context) error {
	resp, err := h.matchUseCase.GetAllMatchResults(context.Background(), &props.GetAllMatchResultsReq{})
	if err != nil {
		return internalServerError
	}

	return c.JSON(http.StatusOK, resp)
}
