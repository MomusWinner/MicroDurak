package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/MommusWinner/MicroDurak/internal/players/v1"
	"github.com/MommusWinner/MicroDurak/services/players"
	"github.com/MommusWinner/MicroDurak/services/players/config"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

func run(ctx context.Context, grpcServer *grpc.Server) error {
	config, err := config.Load()
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, config.DatabaseURL)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	pb.RegisterPlayersServer(grpcServer, players.NewPlayerService(conn, config))

	errChan := make(chan error, 2)

	go func() {
		log.Printf("Starting gRPC server on :%s\n", config.GRPCPort)
		lis, err := net.Listen("tcp", ":"+config.GRPCPort)
		if err != nil {
			errChan <- fmt.Errorf("gRPC listen error: %w", err)
			return
		}

		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quit:
		log.Println("\nShutting down servers...")

		grpcServer.GracefulStop()
		fmt.Println("Servers stopped successfully")
		return nil
	}
}

func main() {
	grpcServer := grpc.NewServer()
	ctx := context.Background()
	if err := run(ctx, grpcServer); err != nil {
		log.Fatal(err)
	}
}
