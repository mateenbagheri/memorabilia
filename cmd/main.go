package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/cluster"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/mateenbagheri/memorabilia/pkg/replication"
	"github.com/mateenbagheri/memorabilia/server"
)

func main() {
	// General flags (in both single node and Raft mode)
	grpcPort := flag.String("port",
		envOrDefault("MEMORABILIA_PORT", "50051"),
		"gRPC server port")

	ttlCleanupMs := flag.Int64("ttl-cleanup-ms",
		envOrDefaultInt64("MEMORABILIA_TTL_CLEANUP_MS", 60000),
		"TTL cleanup job interval in milliseconds")

	// Raft flags (only matters when --node-id is set)
	nodeID := flag.String("node-id",
		envOrDefault("MEMORABILIA_NODE_ID", ""),
		"Unique node ID, e.g. 'n1'. Setting this enables Raft replication mode.")

	raftAddr := flag.String("raft-addr",
		envOrDefault("MEMORABILIA_RAFT_ADDR", "0.0.0.0:7000"),
		"Raft TCP bind address")

	advertiseAddr := flag.String("advertise-addr",
		envOrDefault("MEMORABILIA_ADVERTISE_ADDR", ""),
		"Raft advertise address (defaults to raft-addr; set when behind NAT/Docker)")

	httpMgmtAddr := flag.String("http-mgmt-addr",
		envOrDefault("MEMORABILIA_HTTP_MGMT_ADDR", "0.0.0.0:8081"),
		"HTTP management server address (/raft/join, /raft/leader, /raft/peers)")

	dataDir := flag.String("data-dir",
		envOrDefault("MEMORABILIA_DATA_DIR", "./data"),
		"Base directory for Raft log, stable store, and snapshots (a subdirectory per node-id is created automatically)")

	bootstrap := flag.Bool("bootstrap",
		envOrDefaultBool("MEMORABILIA_BOOTSTRAP", false),
		"Bootstrap a brand-new single-node cluster. Set true on the FIRST node only, on its FIRST run only.")

	leaderHTTP := flag.String("leader-http",
		envOrDefault("MEMORABILIA_LEADER_HTTP", ""),
		"HTTP management address of the cluster leader to join, e.g. '127.0.0.1:8081'")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	repo := core.NewInMemoryCommandRepository()

	// Single node mode
	if *nodeID == "" {
		logger.Info("no --node-id / MEMORABILIA_NODE_ID provided, starting in single-node mode (no replication)")
		srv := server.New(
			server.WithPort(*grpcPort),
			server.WithLogger(logger),
			server.WithCommandsRepository(repo),
			server.WithTTLCleanupTime(*ttlCleanupMs),
		)
		srv.Start()
		return
	}

	// Raft mode
	cfg := &cluster.Config{
		NodeID:         *nodeID,
		RaftBindAddr:   *raftAddr,
		AdvertiseAddr:  *advertiseAddr,
		HTTPMgmtAddr:   *httpMgmtAddr,
		DataDir:        filepath.Join(*dataDir, *nodeID),
		Bootstrap:      *bootstrap,
		LeaderHTTPAddr: *leaderHTTP,
	}

	fsm := replication.NewFSM(repo)

	raftNode, err := replication.NewNode(cfg, fsm, logger)
	if err != nil {
		logger.Error("failed to create raft node", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Non-bootstrap nodes register with the cluster before serving traffic.
	if !cfg.Bootstrap && cfg.LeaderHTTPAddr != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Info("joining cluster via leader", slog.String("leaderHTTP", cfg.LeaderHTTPAddr))
		if err := raftNode.JoinViaLeader(ctx, cfg.LeaderHTTPAddr, cfg.NodeID, cfg.RaftBindAddr); err != nil {
			logger.Error("failed to join cluster", slog.String("error", err.Error()))
			os.Exit(1)
		}
		logger.Info("successfully joined cluster")
	}

	srv := server.New(
		server.WithPort(*grpcPort),
		server.WithLogger(logger),
		server.WithCommandsRepository(repo),
		server.WithRaft(raftNode, fsm, cfg),
		server.WithHTTPMgmtAddr(*httpMgmtAddr),
		server.WithTTLCleanupTime(*ttlCleanupMs),
	)
	srv.Start()
}

// env var helpers
// Important Node: env vars provide defaults, flags override them.

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envOrDefaultBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envOrDefaultInt64(key string, def int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return i
}
