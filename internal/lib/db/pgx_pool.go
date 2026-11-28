package db

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"orglang/go-engine/internal/app"
)

type pgxPoolResource struct {
	pool      *pgxpool.Pool
	closeOnce sync.Once
}

func (r *pgxPoolResource) close() {
	r.closeOnce.Do(r.pool.Close)
}

// Shutdown makes pool ownership explicit for samber/do, including the case
// where the application is stopped before its lifecycle has started.
func (r *pgxPoolResource) Shutdown(context.Context) error {
	r.close()
	return nil
}

// newPoolPgxDo creates a pool whose readiness is owned by the application's
// explicit lifecycle. The resource itself owns closure through do shutdown.
func newPoolPgxDo(props storageProps, lifecycle app.Lifecycle) (*pgxPoolResource, error) {
	config, err := pgxpool.ParseConfig(props.Protocol.Postgres.URL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(props.Driver.Pgx.MaxConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	resource := &pgxPoolResource{pool: pool}

	appendErr := lifecycle.Append(app.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				resource.close()
				return err
			}
			return nil
		},
	})
	if appendErr != nil {
		resource.close()
		return nil, appendErr
	}

	return resource, nil
}
