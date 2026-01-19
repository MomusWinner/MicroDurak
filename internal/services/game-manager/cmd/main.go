package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/core"
	"github.com/MommusWinner/MicroDurak/internal/services/game-manager/delivery/http"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	di := core.NewDi()
	defer core.DisposeCtx(di.Ctx.(*core.Ctx))

	http.AddRoutes(e, di.Handler, di.Ctx)

	errChan := make(chan error, 1)
	go func() {
		di.Ctx.Logger().Info("Game manager server starting...", "port", di.Ctx.Config().GetPort())
		if err := e.Start(":" + di.Ctx.Config().GetPort()); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		di.Ctx.Logger().Error("Server error", "error", err.Error())
		shutdown(e, di)
	case <-quit:
		di.Ctx.Logger().Info("Shutting down game manager service...")
		shutdown(e, di)
		di.Ctx.Logger().Info("Game manager service stopped gracefully")
	}
}

func shutdown(e *echo.Echo, di *core.Di) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		di.Ctx.Logger().Error("Failed to shutdown Echo server", "error", err.Error())
	}
}
