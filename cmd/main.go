package main

import (
	"log"
	"net"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterCommandsServer(s, server.NewCommandServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
