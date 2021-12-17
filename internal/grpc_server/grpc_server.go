package grpc_server

import (
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/history/internal/config"
	"github.com/rusrafkasimov/history/internal/trace"
	"google.golang.org/grpc"
	"net"
)

type HistoryGrpcServer interface {
	RunGrpcServer(regFunc RegisterRPCServer)
	DestroyGrpcServer()
}

type GRPCServerConfig struct {
	host   string
	port   string
	server *grpc.Server
	logger promtail.Client
}

func (s *GRPCServerConfig) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
}

type RegisterRPCServer func(s *GRPCServerConfig)

func NewGRPCServer(config *config.Configuration, logger promtail.Client) *GRPCServerConfig {

	grpcHost, err := config.Get("HISTORY_SERVER_GRPC_HOST")
	if err != nil {
		trace.OnError(logger,nil, err)
		return nil
	}

	grpcPort, err := config.Get("HISTORY_SERVER_GRPC_PORT")
	if err != nil {
		trace.OnError(logger,nil, err)
		return nil
	}

	serverConfig := &GRPCServerConfig{
		host:   grpcHost,
		port:   grpcPort,
		logger: logger,
	}

	var servOpts []grpc.ServerOption
	grpcServer := grpc.NewServer(servOpts...)
	serverConfig.server = grpcServer

	return serverConfig
}

func (s *GRPCServerConfig) RunGrpcServer(regFunc RegisterRPCServer) {

	s.logger.Infof("Run history gRPC server...")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		s.logger.Errorf(err.Error())
	}

	if regFunc == nil {
		s.logger.Errorf("Reg function is nil")
	}

	regFunc(s)

	s.logger.Infof("Listen gRPC at %s:%s", s.host, s.port)

	if err := s.server.Serve(listener); err != nil {
		s.logger.Errorf(err.Error())
	}

	s.logger.Infof("History gRPC server stopped....")
}

func (s *GRPCServerConfig) DestroyGrpcServer() {
	s.logger.Infof("Graceful stop history grpc server...")
	s.server.GracefulStop()
}
