// Disclaimer: Tests in this page are written by LLM
// I have personally taken time to review the generated
// tests here but since this is my personal passion project,
// I prefer any usage of generated content be disclaimed.
package replication

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestFSM(t *testing.T) *FSM {
	t.Helper()
	return NewFSM(core.NewInMemoryCommandRepository())
}

// applyCmd encodes cmd and drives it through the FSM directly, bypassing Raft
// consensus. This lets us test the FSM in isolation.
func applyCmd(t *testing.T, fsm *FSM, cmd *RaftCommand) {
	t.Helper()
	b, err := cmd.Encode()
	require.NoError(t, err, "encoding RaftCommand")

	result := fsm.Apply(&raft.Log{Data: b})
	if result != nil {
		if err, ok := result.(error); ok {
			t.Fatalf("FSM.Apply returned error: %v", err)
		}
	}
}

func TestFSM_Set_And_Get(t *testing.T) {
	fsm := newTestFSM(t)
	ctx := context.Background()

	applyCmd(t, fsm, &RaftCommand{
		Op:         OpSet,
		Key:        "hello",
		Value:      "world",
		Expiration: time.Time{}, // no expiry
	})

	val, err := fsm.Repository().Get(ctx, "hello")
	require.NoError(t, err)
	assert.Equal(t, "world", val)
}

func TestFSM_Set_WithTTL_Expires(t *testing.T) {
	fsm := newTestFSM(t)
	ctx := context.Background()

	applyCmd(t, fsm, &RaftCommand{
		Op:         OpSet,
		Key:        "temp",
		Value:      "42",
		Expiration: time.Now().Add(-1 * time.Second), // already expired
	})

	_, err := fsm.Repository().Get(ctx, "temp")
	assert.ErrorIs(t, err, core.ErrKeyExpiredForGetOp)
}

func TestFSM_Delete(t *testing.T) {
	fsm := newTestFSM(t)
	ctx := context.Background()

	applyCmd(t, fsm, &RaftCommand{Op: OpSet, Key: "k", Value: "v"})
	applyCmd(t, fsm, &RaftCommand{Op: OpDelete, Key: "k"})

	_, err := fsm.Repository().Get(ctx, "k")
	assert.ErrorIs(t, err, core.ErrNotFoundForGetOp)
}

func TestFSM_BatchDelete(t *testing.T) {
	fsm := newTestFSM(t)
	ctx := context.Background()

	applyCmd(t, fsm, &RaftCommand{Op: OpSet, Key: "a", Value: "1"})
	applyCmd(t, fsm, &RaftCommand{Op: OpSet, Key: "b", Value: "2"})
	applyCmd(t, fsm, &RaftCommand{Op: OpSet, Key: "c", Value: "3"})

	result := fsm.Apply(&raft.Log{Data: func() []byte {
		b, _ := (&RaftCommand{Op: OpBatchDelete, Keys: []string{"a", "c"}}).Encode()
		return b
	}()})
	count, ok := result.(int64)
	require.True(t, ok)
	assert.Equal(t, int64(2), count)

	_, err := fsm.Repository().Get(ctx, "a")
	assert.Error(t, err)
	val, err := fsm.Repository().Get(ctx, "b")
	require.NoError(t, err)
	assert.Equal(t, "2", val)
}

func TestFSM_Snapshot_And_Restore(t *testing.T) {
	fsm1 := newTestFSM(t)

	applyCmd(t, fsm1, &RaftCommand{Op: OpSet, Key: "x", Value: "10"})
	applyCmd(t, fsm1, &RaftCommand{Op: OpSet, Key: "y", Value: "20"})

	// Take a snapshot from fsm1
	snap, err := fsm1.Snapshot()
	require.NoError(t, err)

	var buf bytes.Buffer
	sink := &testSnapshotSink{buf: &buf}
	require.NoError(t, snap.Persist(sink))

	// Restore into a fresh fsm2
	fsm2 := newTestFSM(t)
	require.NoError(t, fsm2.Restore(io.NopCloser(&buf)))

	ctx := context.Background()
	for _, tc := range []struct{ k, want string }{{"x", "10"}, {"y", "20"}} {
		got, err := fsm2.Repository().Get(ctx, tc.k)
		require.NoError(t, err, "key %q missing after restore", tc.k)
		assert.Equal(t, tc.want, got)
	}
}

func TestFSM_Apply_UnknownOp_ReturnsError(t *testing.T) {
	fsm := newTestFSM(t)
	b, _ := (&RaftCommand{Op: OpType(99)}).Encode()
	result := fsm.Apply(&raft.Log{Data: b})
	_, isErr := result.(error)
	assert.True(t, isErr, "expected error for unknown op")
}

// testSnapshotSink satisfies raft.SnapshotSink for tests.
type testSnapshotSink struct {
	buf *bytes.Buffer
}

func (s *testSnapshotSink) Write(p []byte) (int, error) { return s.buf.Write(p) }
func (s *testSnapshotSink) Close() error                { return nil }
func (s *testSnapshotSink) ID() string                  { return "test-sink" }
func (s *testSnapshotSink) Cancel() error               { return nil }

var _ raft.SnapshotSink = (*testSnapshotSink)(nil)
