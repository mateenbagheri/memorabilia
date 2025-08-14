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

	commandsRepository core.CommandsRepository
}

func NewCommandServer(
	repo core.CommandsRepository,
) *CommandServer {
	return &CommandServer{
		commandsRepository: repo,
	}
}

func (cs *CommandServer) Echo(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.GetMessage()}, nil
}

func (cs *CommandServer) Get(ctx context.Context, in *api.GetRequest) (*api.GetResponse, error) {
	val, err := cs.commandsRepository.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &api.GetResponse{Value: val}, nil
}

func (cs *CommandServer) Set(ctx context.Context, in *api.SetRequest) (*emptypb.Empty, error) {
	expiration := time.Now().Add(time.Millisecond * time.Duration(in.GetTtl()))
	err := cs.commandsRepository.Set(ctx, in.GetId(), in.GetValue(), expiration)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (cs *CommandServer) Delete(ctx context.Context, in *api.DeleteRequest) (*api.DeleteResponse, error) {
	deleteCount := cs.commandsRepository.Delete(ctx, in.GetId())
	return &api.DeleteResponse{DeleteCount: deleteCount}, nil
}

func (cs *CommandServer) BatchDelete(ctx context.Context, in *api.BatchDeleteRequest) (*api.BatchDeleteResponse, error) {
	deleteCount := cs.commandsRepository.BatchDelete(ctx, in.GetIds())
	return &api.BatchDeleteResponse{DeleteCount: deleteCount}, nil
}

func (cs *CommandServer) GetExpiredKeys(ctx context.Context, in *emptypb.Empty) (*api.GetExpiredKeysResponse, error) {
	expiredKeys, err := cs.commandsRepository.GetExpiredKeys(ctx)
	if err != nil {
		return nil, err
	}
	return &api.GetExpiredKeysResponse{Ids: expiredKeys}, err
}
