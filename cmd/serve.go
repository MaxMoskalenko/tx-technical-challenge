package cmd

import (
	"log"
	"net"
	"os"

	"github.com/MaxMoskalenko/tx-technical-challenge/internal/proto"
	"github.com/MaxMoskalenko/tx-technical-challenge/internal/server"
	"google.golang.org/grpc"
)

func serve() error {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":50051"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAPIServer(grpcServer, server.NewGRPCServer())

	log.Printf("gRPC server listening on %s", addr)
	return grpcServer.Serve(listener)
}
