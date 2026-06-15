package replication

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/raft"
	"github.com/mateenbagheri/memorabilia/pkg/types"
)

type fsmSnapshot struct {
	data map[string]types.ColumnValueWithTTL
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := json.NewEncoder(sink).Encode(f.data); err != nil {
		_ = sink.Cancel()
		return fmt.Errorf("snapshot persist: %w", err)
	}
	return sink.Close()
}

func (f *fsmSnapshot) Release() {}
