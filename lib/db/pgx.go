package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func newOperator(pool *pgxpool.Pool) Operator {
	return &OperatorPgx{pool}
}

func newPgxDriver(dto storageCS, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dto.Protocol.Postgres.URL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(dto.Driver.Pgx.MaxConns)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
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
