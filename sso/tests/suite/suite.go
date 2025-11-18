package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	"sso/internal/config"

	ssov1 "github.com/Aleks12387/gRPC/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T                  // Required to call *testing.T methods within Suite
	Cfg        *config.Config   // Application configuration
	AuthClient ssov1.AuthClient // Client for interacting with gRPC server
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
//
// TODO: for pipeline tests we need to wait for app is ready
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// FIXED: using grpc.NewClient instead of deprecated grpc.DialContext
	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Using insecure connection for tests
	)
	if err != nil {
		t.Fatalf("grpc client creation failed: %v", err)
	}

	// Add connection close to cleanup
	t.Cleanup(func() {
		cc.Close()
	})

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/local.yaml"
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
