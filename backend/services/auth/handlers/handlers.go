package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth/utils"
	"github.com/labstack/echo/v4"
)

func RegisterHandler(c echo.Context) error {
	ac := c.(*AuthContext)
	ctx := c.Request().Context()

	r := new(RegisterRequest)
	if err := c.Bind(r); err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	isEmailTaken, err := ac.DBQueries.CheckEmail(ctx, r.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	if isEmailTaken > 0 {
		return c.String(http.StatusBadRequest, "Email Taken")
	}

	playerId, err := ac.DBQueries.CreatePlayer(ctx, database.CreatePlayerParams{
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

	authId, err := ac.DBQueries.CreateAuth(ctx, database.CreateAuthParams{
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

	jwt, err := utils.GenerateToken(ac.Config.JWTPrivate, authId.String())
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, RegisterResponse{
		Token: jwt,
	})
}

func LoginHandler(c echo.Context) error {
	ac := c.(*AuthContext)

	r := new(LoginRequest)
	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	playerAuth, err := ac.DBQueries.GetAuthByEmail(c.Request().Context(), r.Email)
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

	jwt, err := utils.GenerateToken(ac.Config.JWTPrivate, playerAuth.ID.String())
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: jwt,
	})
}
