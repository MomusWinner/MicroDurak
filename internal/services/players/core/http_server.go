package core

import (
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	app *echo.Echo
	ctx domain.Context
}

type Server interface {
	Start()
	App() *echo.Echo
}

func NewHttpServer(ctx domain.Context) Server {
	app := echo.New()

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
