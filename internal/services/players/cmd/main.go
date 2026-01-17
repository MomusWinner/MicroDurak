package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MommusWinner/MicroDurak/internal/services/players/core"
)

// func run(ctx context.Context, grpcServer *grpc.Server) error {
// 	config, err := config.Load()
// 	if err != nil {
// 		return err
// 	}
//
// 	pool, err := pgxpool.New(ctx, config.DatabaseURL)
// 	if err != nil {
// 		return err
// 	}
// 	defer pool.Close()
//
// 	pb.RegisterPlayersServer(grpcServer, players.NewPlayerService(pool, config))
//
// 	errChan := make(chan error, 2)
//
// 	go func() {
// 		log.Printf("Starting gRPC server on :%s\n", config.GRPCPort)
// 		lis, err := net.Listen("tcp", ":"+config.GRPCPort)
// 		if err != nil {
// 			errChan <- fmt.Errorf("gRPC listen error: %w", err)
// 			return
// 		}
//
// 		if err := grpcServer.Serve(lis); err != nil {
// 			errChan <- fmt.Errorf("gRPC server error: %w", err)
// 		}
// 	}()
//
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//
// 	select {
// 	case err := <-errChan:
// 		return err
// 	case <-quit:
// 		log.Println("\nShutting down servers...")
//
// 		grpcServer.GracefulStop()
// 		fmt.Println("Servers stopped successfully")
// 		return nil
// 	}
// }

func main() {
	var wg sync.WaitGroup

	di := core.NewDi()
	server := core.NewHttpServer(di.Ctx)
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
