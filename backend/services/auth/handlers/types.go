package handlers

import (
	"github.com/MommusWinner/MicroDurak/database"
	"github.com/MommusWinner/MicroDurak/services/auth/config"
	"github.com/labstack/echo/v4"
)

type AuthContext struct {
	echo.Context
	DBQueries *database.Queries
	Config    *config.Config
}

type AuthResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Age      int16  `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
