package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mateenbagheri/memorabilia/pkg/replication"
)

// ScheduleCleanup registers a recurring job that removes expired keys.
//
// The job runs in one of two modes depending on whether Raft is enabled:
//
//   - Raft mode (runReplicatedCleanup): only the current leader scans for
//     expired keys. It then replicates a single BatchDelete command so every
//     node in the cluster removes exactly the same keys at the same position
//     in the log. Followers skip the job entirely — if every node
//     independently scanned and deleted from its own store, clock skew
//     between machines could cause each node to pick a slightly different
//     set of "expired" keys, and the stores would silently diverge.
//
//   - Direct mode (runDirectCleanup): the original single-node behaviour —
//     this node's own store is scanned and pruned locally.
func (s *Server) ScheduleCleanup() {
	cleanUpTimeIntervalInSeconds := s.ttlCleanupTime / 1000
	s.scheduler.ScheduleIntervalJob(fmt.Sprintf("%ds", cleanUpTimeIntervalInSeconds), func() {
		s.logger.Info("running TTL cleanup job",
			slog.Int64("interval_seconds", cleanUpTimeIntervalInSeconds))

		if s.raftNode != nil {
			s.runReplicatedCleanup()
			return
		}
		s.runDirectCleanup()
	})
}

// runReplicatedCleanup scans for expired keys and replicates their deletion
// through Raft. Only the leader performs the scan — see ScheduleCleanup for why.
func (s *Server) runReplicatedCleanup() {
	if !s.raftNode.IsLeader() {
		return
	}

	ctx := context.Background()
	keys, err := s.raftFSM.Repository().GetExpiredKeys(ctx)
	if err != nil {
		s.logger.Error("cleanup: get expired keys failed", slog.String("error", err.Error()))
		return
	}
	if len(keys) == 0 {
		return
	}

	if err := s.raftNode.Apply(&replication.RaftCommand{
		Op:   replication.OpBatchDelete,
		Keys: keys,
	}); err != nil {
		s.logger.Error("cleanup: replicated batch delete failed", slog.String("error", err.Error()))
		return
	}

	s.logger.Info("replicated TTL cleanup", slog.Int("keys_deleted", len(keys)))
}

// runDirectCleanup scans for expired keys and removes them from the local
// store directly. This is the original implementation before raft became a
// thing in the project
func (s *Server) runDirectCleanup() {
	ctx := context.Background()
	deleteCount, err := s.commandsRepository.Cleanup(ctx)
	if err != nil {
		s.logger.Error("cleanup failed", slog.String("error", err.Error()))
		return
	}
	if deleteCount > 0 {
		s.logger.Info("cleaned up expired keys", slog.Int64("keys_deleted", deleteCount))
	}
}
