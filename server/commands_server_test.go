package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/api"
	"github.com/mateenbagheri/memorabilia/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCommandServer_Echo_RandomString(t *testing.T) {
	server := NewCommandServer()

	echoMessage := utils.GenerateRandomString(utils.GenerateRandomNumber(0, 10))
	req := &api.EchoRequest{Message: echoMessage}
	resp, err := server.Echo(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, echoMessage, resp.Message)
}

func TestCommandServer_Echo_NilMessage(t *testing.T) {
	server := NewCommandServer()
	req := &api.EchoRequest{}
	resp, err := server.Echo(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "", resp.Message)
}

func TestEcho_Integration(t *testing.T) {
	// Start a gRPC server in-process
	s, errCh := startTestServer(t)
	defer s.Stop()

	// Create a client connected to the in-process server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("could not create new grpc client: %v", err)
	}
	defer conn.Close()

	c := api.NewCommandsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Generate a random message
	randomMessage := utils.GenerateRandomString(utils.GenerateRandomNumber(0, 100))

	// Test the Echo method
	r, err := c.Echo(ctx, &api.EchoRequest{Message: randomMessage})
	if err != nil {
		t.Fatalf("error while calling echo: %v", err)
	}

	// Assert the result based on message
	assert.Equal(t, randomMessage, r.GetMessage())

	// Check for server errors
	select {
	case serverErr := <-errCh:
		if serverErr != nil {
			t.Fatalf("server error: %v", serverErr)
		}
	default:
	}
}

func startTestServer(t *testing.T) (*grpc.Server, chan error) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterCommandsServer(s, NewCommandServer())

	errCh := make(chan error, 1)
	go func() {
		if err := s.Serve(lis); err != nil {
			errCh <- err
		}
		close(errCh)
	}()
	return s, errCh
}
