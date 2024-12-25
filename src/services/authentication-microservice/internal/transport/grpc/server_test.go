package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"authentication-microservice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func setupTest(t *testing.T) (context.Context, *MockService) {
	mockSrv := new(MockService)
	testLogger := logger.New("test-service")
	ctx := context.WithValue(context.Background(), logger.LoggerKey, testLogger)
	return ctx, mockSrv
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestNew(t *testing.T) {
	ctx, mockSrv := setupTest(t)

	tests := []struct {
		name        string
		grpcPort    int
		restPort    int
		service     Service
		expectError bool
	}{
		{
			name:        "successful server creation",
			grpcPort:    50051,
			restPort:    8080,
			service:     mockSrv,
			expectError: false,
		},
		{
			name:        "nil service",
			grpcPort:    50052,
			restPort:    8081,
			service:     nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := New(ctx, tt.grpcPort, tt.restPort, tt.service)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.NotNil(t, server.grpcServer)
				assert.NotNil(t, server.restServer)
				assert.NotNil(t, server.listener)

				t.Cleanup(func() {
					if server != nil {
						server.grpcServer.Stop()
						if server.listener != nil {
							server.listener.Close()
						}
					}
				})
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	ctx, mockSrv := setupTest(t)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	grpcPort, err := getFreePort()
	assert.NoError(t, err)
	restPort, err := getFreePort()
	assert.NoError(t, err)

	server, err := New(ctx, grpcPort, restPort, mockSrv)
	assert.NoError(t, err)

	errChan := make(chan error)
	go func() {
		errChan <- server.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errChan:
		if err != nil && err != context.Canceled {
			t.Fatalf("server error: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
	}

	err = server.Stop(ctx)
	assert.NoError(t, err)
}

func TestServer_Stop(t *testing.T) {
	ctx, mockSrv := setupTest(t)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	grpcPort, err := getFreePort()
	assert.NoError(t, err)
	restPort, err := getFreePort()
	assert.NoError(t, err)

	server, err := New(ctx, grpcPort, restPort, mockSrv)
	assert.NoError(t, err)

	go func() {
		_ = server.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	err = server.Stop(ctx)
	assert.NoError(t, err)

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err == nil {
		conn.Close()
		t.Error("server should be stopped")
	}
}

func TestContextWithLogger(t *testing.T) {
	testLogger := logger.New("test-service")
	ctx := context.WithValue(context.Background(), logger.LoggerKey, testLogger)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	t.Run("successful case", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			l := logger.GetLoggerFromCtx(ctx)
			assert.NotNil(t, l)
			return "success", nil
		}

		interceptor := ContextWithLogger(testLogger)
		resp, err := interceptor(ctx, "test-request", info, handler)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("error case", func(t *testing.T) {
		expectedErr := fmt.Errorf("test error")
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, expectedErr
		}

		interceptor := ContextWithLogger(testLogger)
		resp, err := interceptor(ctx, "test-request", info, handler)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)
	})
}
