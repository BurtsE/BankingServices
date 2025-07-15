package grpc_server

import (
	"UserService/generated/protobuf"
	"UserService/internal/config"
	"UserService/internal/service"
	"UserService/pkg/middleware"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	protobuf.UnimplementedUserServiceServer
	logger  *logrus.Logger
	service service.UserService
	srv     *grpc.Server
	port    int
}

func NewGrpcServer(logger *logrus.Logger, cfg *config.Config, service service.UserService) *Server {
	return &Server{
		logger:  logger.WithField("server", "grpc").Logger,
		service: service,
		port:    cfg.GrpcPort,
	}
}

func (s *Server) Start() error {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(middleware.NewInterceptorLogger(s.logger)),
			middleware.NewGrpcPanicRecoveryHandler(s.logger),
		))
	protobuf.RegisterUserServiceServer(grpcServer, s)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Fatalf("failed to listen on port %d, error: %v", s.port, err)
	}

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	s.srv = grpcServer

	return grpcServer.Serve(listen)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}
