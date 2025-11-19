package jwt

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwtKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			claims, err := VerifyToken(jwtKey, token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			playerId, err := claims.GetSubject()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims")
			}

			c.Set("playerId", playerId)
			return next(c)
		}
	}
}
