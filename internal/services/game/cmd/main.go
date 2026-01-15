package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/MommusWinner/MicroDurak/internal/game/v1"
	"github.com/MommusWinner/MicroDurak/internal/services/game/config"
	"github.com/MommusWinner/MicroDurak/internal/services/game/controller"
	gameGrpc "github.com/MommusWinner/MicroDurak/internal/services/game/grpc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func run(grpcServer *grpc.Server) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	opt, err := redis.ParseURL(conf.RedisURL)
	if err != nil {
		return err
	}

	client := redis.NewClient(opt)
	channel, err := connectToRabbit(conf)

	if err != nil {
		return err
	}

	errChan := make(chan error, 2)

	gameController := controller.NewGameController(conf, channel, client)
	pb.RegisterGameServer(grpcServer, gameGrpc.NewGameServer(&gameController, conf))

	go func() {
		errChan <- startGrpc(grpcServer, conf)
	}()

	go gameController.ProcessQueues()

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

func startGrpc(grpcServer *grpc.Server, conf *config.Config) error {
	log.Printf("Starting gRPC server on :%s\n", conf.GRPCPort)
	lis, err := net.Listen("tcp", ":"+conf.GRPCPort)
	if err != nil {
		return fmt.Errorf("gRPC listen error: %w", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server error: %w", err)
	}
	return nil
}

func connectToRabbit(conf *config.Config) (*amqp.Channel, error) {
	conn, err := amqp.Dial(conf.RabbitmqURL)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, err
}

func main() {
	grpcServer := grpc.NewServer()

	if err := run(grpcServer); err != nil {
		panic(err)
	}
}
