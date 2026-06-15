package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/cluster"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/mateenbagheri/memorabilia/pkg/replication"
	"github.com/mateenbagheri/memorabilia/pkg/utils/schedule"
	"google.golang.org/grpc"
)

// Server is the top-level process container. It owns the gRPC server, the
// optional HTTP management server (Raft cluster operations), the TTL cleanup
// scheduler, and — when replication is enabled — the Raft node.
//
// Server itself contains no business logic. It is purely a lifecycle:
// construct the dependent servers, start them, wait for a signal, shut
// everything down in order. The actual logic lives in:
//   - CommandServer       (commands_server.go) — gRPC data operations
//   - RaftHTTPHandler     (raft_http.go)        — HTTP cluster management
//   - ScheduleCleanup     (cleanup.go)          — TTL expiry cleanup job
type Server struct {
	ttlCleanupTime     int64 // milliseconds
	logger             *slog.Logger
	grpcPort           string
	httpMgmtAddr       string // e.g. "0.0.0.0:8081"
	grpcServer         *grpc.Server
	httpServer         *http.Server
	scheduler          schedule.CronjobRepository
	commandsRepository core.CommandsRepository

	// Raft fields — nil when running in single-node mode without replication.
	raftNode   *replication.Node
	raftFSM    *replication.FSM
	clusterCfg *cluster.Config
}

// Option configures a Server using the functional-options pattern.
type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) { s.grpcPort = port }
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) { s.logger = logger }
}

func WithTTLCleanupTime(ttlCleanupTime int64) Option {
	return func(s *Server) { s.ttlCleanupTime = ttlCleanupTime }
}

func WithScheduler(scheduler schedule.CronjobRepository) Option {
	return func(s *Server) { s.scheduler = scheduler }
}

func WithCommandsRepository(repo core.CommandsRepository) Option {
	return func(s *Server) { s.commandsRepository = repo }
}

func WithRaft(node *replication.Node, fsm *replication.FSM, cfg *cluster.Config) Option {
	return func(s *Server) {
		s.raftNode = node
		s.raftFSM = fsm
		s.clusterCfg = cfg
	}
}

func WithHTTPMgmtAddr(addr string) Option {
	return func(s *Server) { s.httpMgmtAddr = addr }
}

// New constructs a Server with defaults, then applies options.
func New(options ...Option) *Server {
	s := &Server{
		grpcPort:           "50051",
		httpMgmtAddr:       "0.0.0.0:8081",
		logger:             slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		grpcServer:         grpc.NewServer(),
		scheduler:          schedule.GetRobfigSchedulerInstance(),
		ttlCleanupTime:     60000,
		commandsRepository: core.NewInMemoryCommandRepository(),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", ":"+s.grpcPort)
	if err != nil {
		s.logger.Error("failed to listen", slog.String("error", err.Error()))
		return
	}

	api.RegisterCommandsServer(s.grpcServer, s.buildCommandServer())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	s.startGRPCServer(lis)
	s.startHTTPManagementServer()

	s.ScheduleCleanup()
	s.scheduler.Start()

	<-stop
	s.shutdown()
}

// buildCommandServer returns the gRPC CommandServer in the correct mode:
// Raft-replicated if a node is configured, direct-to-repo otherwise.
func (s *Server) buildCommandServer() *CommandServer {
	if s.raftNode != nil {
		return NewCommandServerWithRaft(s.raftFSM, s.raftNode)
	}
	return NewCommandServer(s.commandsRepository)
}

// startGRPCServer launches the gRPC server on a background goroutine.
func (s *Server) startGRPCServer(lis net.Listener) {
	go func() {
		s.logger.Info("starting gRPC server", slog.String("address", lis.Addr().String()))
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("gRPC server error", slog.String("error", err.Error()))
		}
	}()
}

// startHTTPManagementServer launches the Raft cluster management HTTP server
// (/raft/join, /raft/leader, /raft/peers) on a background goroutine.
// It is a no-op when Raft is not enabled — single-node mode has nothing to
// manage and binds no extra port.
func (s *Server) startHTTPManagementServer() {
	if s.raftNode == nil {
		return
	}

	mux := http.NewServeMux()
	NewRaftHTTPHandler(s.raftNode, s.logger).RegisterRoutes(mux)

	s.httpServer = &http.Server{
		Addr:    s.httpMgmtAddr,
		Handler: mux,
	}

	go func() {
		s.logger.Info("starting HTTP management server", slog.String("address", s.httpMgmtAddr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP management server error", slog.String("error", err.Error()))
		}
	}()
}

// shutdown stops all subsystems in dependency order: refuse new gRPC work
// first, then close the HTTP management server, then stop Raft itself last
func (s *Server) shutdown() {
	s.logger.Info("shutting down...")

	s.grpcServer.GracefulStop()

	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.httpServer.Shutdown(ctx)
	}

	if s.raftNode != nil {
		if err := s.raftNode.Shutdown(); err != nil {
			s.logger.Error("raft shutdown error", slog.String("error", err.Error()))
		}
	}

	time.Sleep(2 * time.Second) // TODO: replace this with channel or waitgroup?
	s.logger.Info("application stopped")
}
