package cluster

type Config struct {
	NodeID string

	// RaftBindAddr is the TCP address Raft listens on for traffic.
	// Example: "0.0.0.0:7000"
	RaftBindAddr string

	// AdvertiseAddr is the address other peers dial to reach this node.
	// Leave empty to use RaftBindAddr. I saw they use this in cases where
	// node sits behind a NAT or Docker service. Since I am not sure how
	// this project will continue, I have decided to add this. No harm done
	// compared to probable benefits
	// Example: "10.0.1.5:7000"
	AdvertiseAddr string

	// HTTPMgmtAddr is the address the HTTP management server listens on.
	// This is separate from the gRPC port and handles /raft/join and /raft/leader.
	// Example: "0.0.0.0:8081"
	HTTPMgmtAddr string

	// DataDir is where BoltDB log/stable stores and snapshots are written.
	// Created on startup if absent. Add to .gitignore.
	// Example: "./data/node1"
	DataDir string

	// Bootstrap must be true only on the very first node of a brand-new cluster.
	// Subsequent nodes set LeaderHTTPAddr and join via HTTP.
	Bootstrap bool

	// LeaderHTTPAddr is the HTTP management address of the current leader.
	// Non-bootstrap nodes send their JoinRequest here at startup.
	// Example: "10.0.1.5:8081"
	LeaderHTTPAddr string
}
