package server

import (
	"context"

	"github.com/mateenbagheri/memorabilia/api"
)

type CommandServer struct {
	api.UnimplementedCommandsServer
}

func NewCommandServer() *CommandServer {
	return &CommandServer{}
}

func (s *CommandServer) Echo(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.Message}, nil
}
