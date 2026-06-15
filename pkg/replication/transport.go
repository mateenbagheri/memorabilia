package replication

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/raft"
)

const (
	raftTCPTimeout = 5 * time.Second
	raftMaxPool    = 5
)

func NewTCPTransport(bindAddr, advertiseAddr string) (*raft.NetworkTransport, error) {
	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		return nil, fmt.Errorf("replication: resolve bind addr %q: %w", bindAddr, err)
	}

	var advertise net.Addr
	if advertiseAddr != "" && advertiseAddr != bindAddr {
		advertise, err = net.ResolveTCPAddr("tcp", advertiseAddr)
		if err != nil {
			return nil, fmt.Errorf("replication: resolve advertise addr %q: %w", advertiseAddr, err)
		}
	}

	transport, err := raft.NewTCPTransportWithLogger(
		addr.String(),
		advertise,
		raftMaxPool,
		raftTCPTimeout,
		nil, // TODO: figure out if I will change this to my slog or not later on
	)
	if err != nil {
		return nil, fmt.Errorf("replication: create TCP transport: %w", err)
	}

	return transport, nil
}
