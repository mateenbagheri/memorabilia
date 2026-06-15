package cluster

import (
	"fmt"

	"github.com/hashicorp/raft"
)

type RaftNode interface {
	Join(nodeID, raftAddr string) error
	IsLeader() bool
	LeaderRaftAddr() string
	Raft() *raft.Raft
}

type Membership struct {
	node RaftNode
}

func NewMembership(node RaftNode) *Membership {
	return &Membership{node: node}
}

// AddVoter registers a new voting member. Only the leader can do this.
func (m *Membership) AddVoter(nodeID, raftAddr string) error {
	if !m.node.IsLeader() {
		return fmt.Errorf("membership: not the leader (leader raft addr: %q)", m.node.LeaderRaftAddr())
	}
	return m.node.Join(nodeID, raftAddr)
}

// RemoveServer removes a node from the cluster by ID. Obviously only the leader can do this.
func (m *Membership) RemoveServer(nodeID string) error {
	if !m.node.IsLeader() {
		return fmt.Errorf("membership: not the leader")
	}
	f := m.node.Raft().RemoveServer(raft.ServerID(nodeID), 0, 0)
	if err := f.Error(); err != nil {
		return fmt.Errorf("membership: remove server %q: %w", nodeID, err)
	}
	return nil
}

// Servers returns the current cluster all member list
func (m *Membership) Servers() ([]raft.Server, error) {
	f := m.node.Raft().GetConfiguration()
	if err := f.Error(); err != nil {
		return nil, fmt.Errorf("membership: get configuration: %w", err)
	}
	return f.Configuration().Servers, nil
}

// LeaderRaftAddr returns the Raft transport address of the current leader known to this node.
func (m *Membership) LeaderRaftAddr() string {
	return m.node.LeaderRaftAddr()
}
