//go:build integration

package cases

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"telephone-service/internal/tests/containers"

	"github.com/testcontainers/testcontainers-go"
)

var databaseURL string

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) (code int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	postgres, err := containers.StartPostgres(ctx)
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := testcontainers.TerminateContainer(postgres, testcontainers.StopContext(cleanupCtx)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			code = 1
		}
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	databaseURL, err = postgres.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return m.Run()
}
