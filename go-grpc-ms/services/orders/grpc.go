package main

import (
	"log"
	"net"

	handler "github.com/thutasann/go-grpc-ms/services/orders/handler/orders"
	"github.com/thutasann/go-grpc-ms/services/orders/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// register the grpc services
	ordserService := service.NewOrderService()
	handler.NewGrpcOrdersService(grpcServer, ordserService)

	log.Println("Starting gRPC server on ", s.addr)

	return grpcServer.Serve(lis)
}
