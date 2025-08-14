package benchmarks

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/core"
	"github.com/mateenbagheri/memorabilia/pkg/utils/testutil"
	"github.com/mateenbagheri/memorabilia/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func BenchmarkEcho(b *testing.B) {
	// Start a gRPC server in-process
	s, errCh := startBenchmarkServer(b)
	defer s.Stop()

	// Create a client connected to the in-process server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatalf("could not create new grpc client: %v", err)
	}
	defer conn.Close()

	c := api.NewCommandsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Benchmark the Echo method
	for i := 0; i < b.N; i++ {
		randomMessage := testutil.GenerateRandomString(testutil.GenerateRandomInteger(0, 100))
		_, err := c.Echo(ctx, &api.EchoRequest{Message: randomMessage})
		if err != nil {
			b.Fatalf("error while calling echo: %v", err)
		}
	}

	// Check for server errors
	select {
	case serverErr := <-errCh:
		if serverErr != nil {
			b.Fatalf("server error: %v", serverErr)
		}
	default:
	}
}

func startBenchmarkServer(b *testing.B) (*grpc.Server, chan error) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		b.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	repo := core.NewInMemoryCommandRepository()
	api.RegisterCommandsServer(s, server.NewCommandServer(repo))

	errCh := make(chan error, 1)
	go func() {
		if err := s.Serve(lis); err != nil {
			errCh <- err
		}
		close(errCh)
	}()
	return s, errCh
}
