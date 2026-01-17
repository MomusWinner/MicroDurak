package core

import (
	"net"

	pb "github.com/MommusWinner/MicroDurak/internal/contracts/players/v1"
	grpc2 "github.com/MommusWinner/MicroDurak/internal/services/players/delivery/grpc"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain"
	"github.com/MommusWinner/MicroDurak/internal/services/players/domain/cases"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGrpcServer(ctx domain.Context, playerUseCase *cases.PlayerUseCase, matchUseCase *cases.MatchUseCase) *GrpcServer {
	grpcListener, err := net.Listen("tcp", ":"+ctx.Config().GetGRPCPort())
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	service := grpc2.NewPlayerService(ctx, playerUseCase, matchUseCase)
	pb.RegisterPlayersServer(grpcServer, service)

	return &GrpcServer{
		listener:   grpcListener,
		grpcServer: grpcServer,
	}
}

func (s *GrpcServer) Start() {
	if err := s.grpcServer.Serve(s.listener); err != nil {
		panic("failed to serve: " + err.Error())
	}
}
