package main

import (
	"github.com/MommusWinner/MicroDurak/internal/services/auth/api/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/core"
	"github.com/MommusWinner/MicroDurak/lib/validate"
	"github.com/go-playground/validator"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = validate.NewHttpValidator(validator.New())

	di := core.NewDi()

	v1.AddRoutes(e, di.AuthHandler)

	err := e.Start(":" + di.Ctx.Config().GetPort())
	if err != nil {
		panic(err)
	}
}
