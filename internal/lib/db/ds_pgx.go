package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
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
