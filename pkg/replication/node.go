package replication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/mateenbagheri/memorabilia/pkg/cluster"
)

const (
	applyTimeout   = 5 * time.Second
	snapshotRetain = 3
)

type Node struct {
	raft      *raft.Raft
	transport *raft.NetworkTransport
	cfg       *cluster.Config
	logger    *slog.Logger
}

// NewNode simply creates, configures, and starts a Raft node.
func NewNode(cfg *cluster.Config, fsm *FSM, logger *slog.Logger) (*Node, error) {
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return nil, fmt.Errorf("node: mkdir %q: %w", cfg.DataDir, err)
	}

	transport, err := NewTCPTransport(cfg.RaftBindAddr, cfg.AdvertiseAddr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.DataDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(cfg.DataDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(cfg.DataDir, snapshotRetain, os.Stderr)
	if err != nil {
		return nil, err
	}

	raftCfg := raft.DefaultConfig()
	raftCfg.LocalID = raft.ServerID(cfg.NodeID)

	r, err := raft.NewRaft(raftCfg, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("node: new raft: %w", err)
	}

	if cfg.Bootstrap {
		hasState, err := raft.HasExistingState(logStore, stableStore, snapshotStore)
		if err != nil {
			return nil, fmt.Errorf("node: check existing state: %w", err)
		}
		if !hasState {
			bootCfg := raft.Configuration{
				Servers: []raft.Server{
					{ID: raftCfg.LocalID, Address: transport.LocalAddr()},
				},
			}
			if f := r.BootstrapCluster(bootCfg); f.Error() != nil {
				return nil, fmt.Errorf("node: bootstrap: %w", f.Error())
			}
			logger.Info("bootstrapped new Raft cluster",
				slog.String("nodeID", cfg.NodeID),
				slog.String("raftAddr", cfg.RaftBindAddr),
			)
		}
	}

	return &Node{raft: r, transport: transport, cfg: cfg, logger: logger}, nil
}

func (n *Node) IsLeader() bool {
	return n.raft.State() == raft.Leader
}

func (n *Node) LeaderRaftAddr() string {
	addr, _ := n.raft.LeaderWithID()
	return string(addr)
}

func (n *Node) Raft() *raft.Raft {
	return n.raft
}

func (n *Node) Apply(cmd *RaftCommand) error {
	b, err := cmd.Encode()
	if err != nil {
		return fmt.Errorf("node apply: encode: %w", err)
	}

	f := n.raft.Apply(b, applyTimeout)
	if err := f.Error(); err != nil {
		return fmt.Errorf("node apply: raft: %w", err)
	}

	// TODO: check if this can be bettered
	if resp, ok := f.Response().(error); ok && resp != nil {
		return fmt.Errorf("node apply: fsm: %w", resp)
	}

	return nil
}

type JoinRequest struct {
	NodeID   string `json:"node_id"`
	RaftAddr string `json:"raft_addr"`
}

// Join adds nodeID/raftAddr as a new voting member. Impo: Must be called on the leader.
func (n *Node) Join(nodeID, raftAddr string) error {
	n.logger.Info("adding voter", slog.String("nodeID", nodeID), slog.String("raftAddr", raftAddr))
	f := n.raft.AddVoter(
		raft.ServerID(nodeID),
		raft.ServerAddress(raftAddr),
		0, 0,
	)
	if err := f.Error(); err != nil {
		return fmt.Errorf("node join: AddVoter %q: %w", nodeID, err)
	}
	return nil
}

func (n *Node) JoinViaLeader(ctx context.Context, leaderHTTPAddr, nodeID, raftAddr string) error {
	body, err := json.Marshal(JoinRequest{NodeID: nodeID, RaftAddr: raftAddr})
	if err != nil {
		return fmt.Errorf("node joinViaLeader: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://"+leaderHTTPAddr+"/raft/join",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("node joinViaLeader: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("node joinViaLeader: POST %q: %w", leaderHTTPAddr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node joinViaLeader: leader rejected join with status %d", resp.StatusCode)
	}
	return nil
}

func (n *Node) Shutdown() error {
	if f := n.raft.Shutdown(); f.Error() != nil {
		return fmt.Errorf("node shutdown raft: %w", f.Error())
	}
	if err := n.transport.Close(); err != nil {
		return fmt.Errorf("node shutdown transport: %w", err)
	}
	return nil
}
