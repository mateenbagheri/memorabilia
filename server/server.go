package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/mateenbagheri/memorabilia/pkg/utils/schedule"
	"google.golang.org/grpc"
)

type Server struct {
	ttlCleanupTime     int64 // In Milliseconds
	logger             *slog.Logger
	port               string
	grpcServer         *grpc.Server
	scheduler          schedule.CronjobRepository
	commandsRepository core.CommandsRepository
}

type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithTTLCleanupTime(ttlCleanupTime int64) Option {
	return func(s *Server) {
		s.ttlCleanupTime = ttlCleanupTime
	}
}

func WithScheduler(scheduler schedule.CronjobRepository) Option {
	return func(s *Server) {
		s.scheduler = scheduler
	}
}

func WithCommandsRepository(repo core.CommandsRepository) Option {
	return func(s *Server) {
		s.commandsRepository = repo
	}
}

func New(options ...Option) *Server {
	s := &Server{
		port:               "50051",
		logger:             slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		grpcServer:         grpc.NewServer(),
		scheduler:          schedule.GetRobfigSchedulerInstance(),
		ttlCleanupTime:     4000, // TODO: tweak this later on also add env set to this.
		commandsRepository: core.NewInMemoryCommandRepository(),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.logger.Error("failed to listen", slog.String("error", err.Error()))
		return
	}

	commandServer := NewCommandServer(s.commandsRepository)
	api.RegisterCommandsServer(s.grpcServer, commandServer)

	// Channel to catch OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		s.logger.Info("Starting gRPC server", slog.String("address", lis.Addr().String()))
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve", slog.String("error", err.Error()))
		}
	}()

	s.ScheduleCleanup()
	s.scheduler.Start()

	<-stop
	s.logger.Info("Shutting down server...")
	s.grpcServer.GracefulStop()
	time.Sleep(2 * time.Second)
	s.logger.Info("Application stopped")
}

func (s *Server) ScheduleCleanup() {
	cleanUpTimeIntervalInSeconds := s.ttlCleanupTime / 1000
	s.scheduler.ScheduleIntervalJob(fmt.Sprintf("%ds", cleanUpTimeIntervalInSeconds), func() {
		s.logger.Info("Running TTL cleanup job ...", slog.Int64("interval in seconds", cleanUpTimeIntervalInSeconds))

		// Create a dedicated new context for the cleanup task
		// Note: This is a dedicated ctx because ScheduleCleanup is
		// a background task not driven from a user/client request
		ctx := context.Background()

		deleteCount, err := s.commandsRepository.Cleanup(ctx)
		if err != nil {
			s.logger.Error("ScheduleCleanup is failing", slog.String("error", err.Error()))
		}

		if deleteCount > 0 {
			s.logger.Info("ScheduleCleanup cleaned up keys from memory", slog.Int64("keys deleted", deleteCount))
		}
	})
}
