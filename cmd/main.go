package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/utils/schedule"
	"github.com/mateenbagheri/memorabilia/server"
	"google.golang.org/grpc"
)

func main() {
	// Set up slog for structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("failed to listen", slog.String("error", err.Error()))
		return
	}

	scheduler := schedule.GetRobfigSchedulerInstance()
	scheduler.Start()
	defer scheduler.Stop() // Ensure the scheduler stops on exit

	s := grpc.NewServer()
	api.RegisterCommandsServer(s, server.NewCommandServer())

	// Channel to catch OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run gRPC server in a goroutine
	go func() {
		slog.Info("Starting gRPC server", slog.String("address", lis.Addr().String()))
		if err := s.Serve(lis); err != nil {
			slog.Error("failed to serve", slog.String("error", err.Error()))
		}
	}()

	// Block until a signal is received
	<-stop
	slog.Info("Shutting down server...")

	// Gracefully stop the gRPC server
	s.GracefulStop()

	// Allow time for cleanup
	time.Sleep(2 * time.Second) // Adjust as needed
	slog.Info("Application stopped")
}
