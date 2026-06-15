package replication

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/raft"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/mateenbagheri/memorabilia/pkg/types"
)

type FSM struct {
	repo core.CommandsRepository
}

func NewFSM(repo core.CommandsRepository) *FSM {
	return &FSM{repo: repo}
}

func (fsm *FSM) Repository() core.CommandsRepository {
	return fsm.repo
}

func (fsm *FSM) Apply(l *raft.Log) any {
	cmd, err := DecodeCommand(l.Data)
	if err != nil {
		return fmt.Errorf("fsm apply: decode: %w", err)
	}

	ctx := context.Background()

	switch cmd.Op {
	case OpSet:
		return fsm.repo.Set(ctx, cmd.Key, cmd.Value, cmd.Expiration)
	case OpDelete:
		count := fsm.repo.Delete(ctx, cmd.Key)
		return count
	case OpBatchDelete:
		count := fsm.repo.BatchDelete(ctx, cmd.Keys)
		return count
	default:
		return fmt.Errorf("fsm apply: unknown op %d", cmd.Op)
	}
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	data, err := fsm.repo.Dump()
	if err != nil {
		return nil, fmt.Errorf("fsm snapshot: %w", err)
	}

	return &fsmSnapshot{data: data}, nil
}

func (fsm *FSM) Restore(rc io.ReadCloser) error {
	defer rc.Close()

	var data map[string]types.ColumnValueWithTTL
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		return fmt.Errorf("fsm restore: decode: %w", err)
	}

	return fsm.repo.Load(data)
}
