package http

import (
	"errors"
	"net/http"

	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/cases"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/props"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PlayerHandler struct {
	useCase *cases.PlayerUseCase
}

func NewPlayerHandler(useCase *cases.PlayerUseCase) *PlayerHandler {
	return &PlayerHandler{
		useCase: useCase,
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
	resp, err := h.useCase.GetAll(props.GetAllPlayersReq{})

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

	resp, err := h.useCase.GetById(props.GetPlayerByIdReq{Id: playerId})

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
