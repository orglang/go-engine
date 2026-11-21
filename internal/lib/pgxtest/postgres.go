package pgxtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewPostgresPool starts an isolated PostgreSQL container and returns a pool
// limited to one connection so temporary test schemas remain visible.
func NewPostgresPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.11-alpine3.24",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "entpgx_test",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(time.Minute),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Minute)
		defer cleanupCancel()
		if err := container.Terminate(cleanupCtx); err != nil {
			t.Errorf("terminate PostgreSQL container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("get PostgreSQL container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("get mapped PostgreSQL port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/entpgx_test?sslmode=disable", host, port.Port())
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse PostgreSQL connection string: %v", err)
	}
	config.MaxConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("open PostgreSQL pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping PostgreSQL container: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Minute)
		defer cleanupCancel()
		if err := container.Terminate(cleanupCtx); err != nil {
			t.Errorf("terminate PostgreSQL container: %v", err)
		}
	})
	t.Cleanup(pool.Close)

	return pool
}
