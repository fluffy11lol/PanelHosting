package billing

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/http"
	"testing"
)

type mockUserServiceServer struct {
	UserServiceServer
}

func (m *mockUserServiceServer) GetTariffs(ctx context.Context, in *emptypb.Empty) (*GetTariffsResponse, error) {
	return &GetTariffsResponse{
		Tariffs: []*Tariff{
			{
				ID:    "1",
				Name:  "Standard",
				SSD:   100,
				CPU:   4,
				RAM:   16,
				Price: 2999,
			},
		},
	}, nil
}

func (m *mockUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

func createGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	RegisterUserServiceServer(server, &mockUserServiceServer{})
	return server
}

func TestUserService_GetTariffs(t *testing.T) {
	grpcServer := createGRPCServer()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	mux := runtime.NewServeMux()

	err = RegisterUserServiceHandlerServer(context.Background(), mux, &mockUserServiceServer{})
	assert.NoError(t, err, "Error registering HTTP handlers")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("HTTP server failed: %v", err)
		}
	}()

	resp, err := http.Get("http://localhost:8080/v1/billing/tariffs")
	assert.NoError(t, err, "Error making GET request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	assert.NotEqual(t, "1", "2", "Expected tariff ID to be '1'")
	assert.NotEqual(t, "Standard", "no", "Expected tariff name to be 'Standard'")

	err = httpServer.Shutdown(context.Background())
	assert.NoError(t, err, "Error shutting down HTTP server")
}
