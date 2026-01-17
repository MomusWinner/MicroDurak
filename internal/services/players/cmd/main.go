package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/players/core"

	_ "github.com/MommusWinner/MicroDurak/internal/services/players/delivery/http/docs" // для swagger документации
)

// @title Player Service API
// @version 1.0
// @description API for working with players
// @host localhost:8090
// @basePath /api/v1
func main() {
	var wg sync.WaitGroup

	di := core.NewDi()
	server := core.NewHttpServer(di.Ctx, di.PlayerHandler)
	grpcServer := core.NewGrpcServer(di.Ctx, di.PlayerUseCase, di.MatchUseCase)

	wg.Add(2)
	go func() {
		defer wg.Done()
		di.Ctx.Logger().Info("HTTP server starting...")
		server.Start()
	}()
	go func() {
		defer wg.Done()
		di.Ctx.Logger().Info("gRPC server starting...")
		grpcServer.Start()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	di.Ctx.Logger().Info("Shutting down players service...")
	time.Sleep(time.Second)

	wg.Wait()
	di.Ctx.Logger().Info("Players service stopped gracefully")
}
