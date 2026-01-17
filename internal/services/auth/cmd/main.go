package main

import (
	"github.com/MommusWinner/MicroDurak/internal/services/auth/core"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/delivery/http"
	"github.com/MommusWinner/MicroDurak/lib/validate"
	"github.com/go-playground/validator"
	"github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"

	_ "github.com/MommusWinner/MicroDurak/internal/services/auth/delivery/http/docs" // для swagger документации
)

// @title Auth Service API
// @version 1.0
// @description API for authentication and user registration
// @host localhost:8080
// @basePath /api/v1/auth
func main() {
	e := echo.New()
	e.Validator = validate.NewHttpValidator(validator.New())

	di := core.NewDi()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	http.AddRoutes(e, di.AuthHandler)

	err := e.Start(":" + di.Ctx.Config().GetPort())
	if err != nil {
		panic(err)
	}
}
