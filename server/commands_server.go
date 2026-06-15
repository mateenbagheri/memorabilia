package server

import (
	"context"
	"time"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/mateenbagheri/memorabilia/pkg/replication"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	applyTimeout   = 5 * time.Second
	snapshotRetain = 3
)

type CommandServer struct {
	api.UnimplementedCommandsServer

	// repo is used for all reads and for writes in direct mode.
	// In Raft mode it points at the same store that the FSM writes to,
	// so reads always see the latest committed state on this node.
	repo core.CommandsRepository
	fsm  *replication.FSM
	node *replication.Node
}

func NewCommandServer(
	repo core.CommandsRepository,
) *CommandServer {
	return &CommandServer{repo: repo}
}

func NewCommandServerWithRaft(
	fsm *replication.FSM,
	node *replication.Node,
) *CommandServer {
	return &CommandServer{
		fsm:  fsm,
		node: node,
		repo: fsm.Repository(),
	}
}

func (cs *CommandServer) isRaftMode() bool {
	return cs.node != nil
}

// -- Handlers --

func (cs *CommandServer) Echo(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.GetMessage()}, nil
}

func (cs *CommandServer) Get(ctx context.Context, in *api.GetRequest) (*api.GetResponse, error) {
	val, err := cs.repo.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &api.GetResponse{Value: val}, nil
}

func (cs *CommandServer) Set(ctx context.Context, in *api.SetRequest) (*emptypb.Empty, error) {
	ttl := time.Duration(in.Ttl)
	expiration := time.Now().Add(ttl * time.Millisecond)
	if cs.isRaftMode() {
		if err := cs.requireleader(); err != nil {
			return nil, err
		}

		command := &replication.RaftCommand{
			Op:    replication.OpSet,
			Key:   in.GetId(),
			Value: in.GetValue(),
		}

		if ttl > 0 {
			command.Expiration = expiration
		}

		if err := cs.node.Apply(command); err != nil {
			return nil, status.Errorf(codes.Internal, "set (raft): %v", err)
		}
		return &emptypb.Empty{}, nil
	}

	err := cs.repo.Set(ctx, in.GetId(), in.GetValue(), expiration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "set: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (cs *CommandServer) Delete(ctx context.Context, in *api.DeleteRequest) (*api.DeleteResponse, error) {
	if cs.isRaftMode() {
		if err := cs.requireleader(); err != nil {
			return nil, err
		}

		if err := cs.node.Apply(&replication.RaftCommand{
			Op:  replication.OpDelete,
			Key: in.GetId(),
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "delete (raft): %v", err)
		}
		return &api.DeleteResponse{DeleteCount: 1}, nil
	}
	deleteCount := cs.repo.Delete(ctx, in.GetId())
	return &api.DeleteResponse{DeleteCount: deleteCount}, nil
}

func (cs *CommandServer) BatchDelete(ctx context.Context, in *api.BatchDeleteRequest) (*api.BatchDeleteResponse, error) {
	if cs.isRaftMode() {
		if err := cs.requireleader(); err != nil {
			return nil, err
		}
		if err := cs.node.Apply(&replication.RaftCommand{
			Op:   replication.OpBatchDelete,
			Keys: in.GetIds(),
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "batch delete (raft): %v", err)
		}
		return &api.BatchDeleteResponse{DeleteCount: int64(len(in.GetIds()))}, nil
	}
	deleteCount := cs.fsm.Repository().BatchDelete(ctx, in.GetIds())
	return &api.BatchDeleteResponse{DeleteCount: deleteCount}, nil
}

func (cs *CommandServer) GetExpiredKeys(ctx context.Context, in *emptypb.Empty) (*api.GetExpiredKeysResponse, error) {
	expiredKeys, err := cs.fsm.Repository().GetExpiredKeys(ctx)
	if err != nil {
		return nil, err
	}
	return &api.GetExpiredKeysResponse{Ids: expiredKeys}, err
}

// requireLeader returns a gRPC FailedPrecondition error when this node is not
// the leader. The error message includes the leader's Raft address so clients
// can locate the leader and retry.
func (cs *CommandServer) requireleader() error {
	if cs.node.IsLeader() {
		return nil
	}
	leader := cs.node.LeaderRaftAddr()
	if leader == "" {
		return status.Error(codes.Unavailable, "no leader elected yet, retry shortly")
	}
	return status.Errorf(codes.FailedPrecondition,
		"not the leader; current leader raft addr is %q", leader)
}
