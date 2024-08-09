package server

import (
	"context"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommandServer struct {
	api.UnimplementedCommandsServer
}

var inMemoryCommandRepository = core.NewInMemoryCommandRepository()

func NewCommandServer() *CommandServer {
	return &CommandServer{}
}

func (s *CommandServer) Echo(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.Message}, nil
}

func (s *CommandServer) Get(ctx context.Context, in *api.GetRequest) (*api.GetResponse, error) {
	val, err := inMemoryCommandRepository.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &api.GetResponse{Value: val}, nil
}

func (s *CommandServer) Set(ctx context.Context, in *api.SetRequest) (*emptypb.Empty, error) {
	err := inMemoryCommandRepository.Set(ctx, in.Id, in.Value)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
