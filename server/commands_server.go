package server

import (
	"context"
	"time"

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
	return &api.EchoResponse{Message: in.GetMessage()}, nil
}

func (s *CommandServer) Get(ctx context.Context, in *api.GetRequest) (*api.GetResponse, error) {
	val, err := inMemoryCommandRepository.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &api.GetResponse{Value: val}, nil
}

func (s *CommandServer) Set(ctx context.Context, in *api.SetRequest) (*emptypb.Empty, error) {
	expiration := time.Now().Add(time.Millisecond * time.Duration(in.GetTtl()))
	err := inMemoryCommandRepository.Set(ctx, in.GetId(), in.GetValue(), expiration)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CommandServer) Delete(ctx context.Context, in *api.DeleteRequest) (*api.DeleteResponse, error) {
	deleteCount := inMemoryCommandRepository.Delete(ctx, in.GetId())
	return &api.DeleteResponse{DeleteCount: deleteCount}, nil
}

func (s *CommandServer) BatchDelete(ctx context.Context, in *api.BatchDeleteRequest) (*api.BatchDeleteResponse, error) {
	deleteCount := inMemoryCommandRepository.BatchDelete(ctx, in.GetIds())
	return &api.BatchDeleteResponse{DeleteCount: deleteCount}, nil
}
