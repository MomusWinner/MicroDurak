package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/MommusWinner/MicroDurak/internal/database"
	"github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/MommusWinner/MicroDurak/services/auth/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DBQueries     *database.Queries
	Config        *config.Config
	PlayersClient players.PlayersClient
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	type RegisterRequest struct {
		Name     string `json:"name" validate:"required"`
		Age      int16  `json:"age" validate:"required,gte=0,lte=130"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	r := new(RegisterRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	isEmailTaken, err := h.DBQueries.CheckEmail(ctx, r.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	if isEmailTaken > 0 {
		return c.String(http.StatusBadRequest, "Email Taken")
	}

	rep, err := h.PlayersClient.CreatePlayer(ctx, &players.CreatePlayerRequest{
		Name: r.Name,
		Age:  int32(r.Age),
	})
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	hashedPassword, err := utils.HashPassword(r.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	playerId, err := uuid.Parse(rep.Id)
	authId, err := h.DBQueries.CreateAuth(ctx, database.CreateAuthParams{
		PlayerID: pgtype.UUID{Valid: true, Bytes: playerId},
		Email:    r.Email,
		Password: hashedPassword,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	jwt, err := utils.GenerateToken(h.Config.JWTPrivate, authId.String())
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token: jwt,
	})
}

func (h *Handler) Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	r := new(LoginRequest)
	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	playerAuth, err := h.DBQueries.GetAuthByEmail(c.Request().Context(), r.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.String(http.StatusForbidden, "Login failed")
		} else {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	if !utils.CheckPasswordHash(r.Password, playerAuth.Password) {
		return c.String(http.StatusForbidden, "Login failed")
	}

	jwt, err := utils.GenerateToken(h.Config.JWTPrivate, playerAuth.ID.String())
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: jwt,
	})
}
