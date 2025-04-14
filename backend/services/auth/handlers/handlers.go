package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/MommusWinner/MicroDurak/services/auth/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DBQueries *database.Queries
	Config    *config.Config
}

func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	r := new(RegisterRequest)
	if err := c.Bind(r); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	isEmailTaken, err := h.DBQueries.CheckEmail(ctx, r.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	if isEmailTaken > 0 {
		return c.String(http.StatusBadRequest, "Email Taken")
	}

	playerId, err := h.DBQueries.CreatePlayer(ctx, database.CreatePlayerParams{
		Name: r.Name,
		Age:  r.Age,
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

	authId, err := h.DBQueries.CreateAuth(ctx, database.CreateAuthParams{
		PlayerID: playerId,
		Email:    r.Email,
		Password: hashedPassword,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	type RegisterResponse struct {
		Token string `json:"token"`
	}

	jwt, err := utils.GenerateToken(h.Config.JWTPrivate, authId.String())
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, RegisterResponse{
		Token: jwt,
	})
}

func (h *Handler) Login(c echo.Context) error {
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
