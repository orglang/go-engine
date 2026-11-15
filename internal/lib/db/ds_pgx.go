package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"orglang/go-engine/internal/app"
)

func newPoolPgx(props storageProps, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	config, parseErr := pgxpool.ParseConfig(props.Protocol.Postgres.URL)
	if parseErr != nil {
		return nil, parseErr
	}
	config.MaxConns = int32(props.Driver.Pgx.MaxConns)
	pool, createErr := pgxpool.NewWithConfig(context.Background(), config)
	if createErr != nil {
		return nil, createErr
	}
	lc.Append(
		fx.Hook{
			OnStart: pool.Ping,
			OnStop: func(ctx context.Context) error {
				go pool.Close()
				return nil
			},
		},
	)
	return pool, nil
}

// newPoolPgxDo creates a pool whose readiness and shutdown are owned by the
// application's explicit lifecycle. A failed readiness check closes the pool
// immediately because the lifecycle does not stop a hook whose start failed.
func newPoolPgxDo(props storageProps, lifecycle app.Lifecycle) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(props.Protocol.Postgres.URL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(props.Driver.Pgx.MaxConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := lifecycle.Append(app.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				pool.Close()
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	}); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
