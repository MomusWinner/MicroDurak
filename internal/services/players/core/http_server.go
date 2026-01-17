package core

import (
	"github.com/MommusWinner/MicroDurak/internal/services/players/delivery/http"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/lib/validate"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

type HttpServer struct {
	app *echo.Echo
	ctx domain.Context
}

type Server interface {
	Start()
	App() *echo.Echo
}

func NewHttpServer(ctx domain.Context, playerHandler *http.PlayerHandler) Server {
	app := echo.New()
	app.Validator = validate.NewHttpValidator(validator.New())

	app.GET("/swagger/*", echoSwagger.WrapHandler)

	http.AddRoutes(app, playerHandler)

	return &HttpServer{
		app: app,
		ctx: ctx,
	}
}

func (s *HttpServer) Start() {
	err := s.app.Start(":" + s.ctx.Config().GetHTTPPort())
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic("http server inst start successfully")
	}
}

func (s *HttpServer) App() *echo.Echo {
	return s.app
}
