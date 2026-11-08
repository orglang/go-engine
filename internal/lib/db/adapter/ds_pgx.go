// Package adapter provides the pgx implementation of the transaction contract.
package adapter

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"orglang/go-engine/internal/lib/tx"
)

type pgxPool interface {
	Begin(context.Context) (pgx.Tx, error)
	Acquire(context.Context) (*pgxpool.Conn, error)
}

type transactorPgx struct {
	pool pgxPool
}

// NewTransactorPgx creates a PostgreSQL-backed transactor.
func NewTransactorPgx(pool *pgxpool.Pool) tx.Transactor {
	return &transactorPgx{pool: pool}
}

func (t *transactorPgx) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) (retErr error) {
	if fn == nil {
		return errors.New("transaction callback is nil")
	}

	pgxTx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	// Rollback is attempted on every exit. ErrTxClosed is expected after a
	// successful commit; any other rollback error is joined with the result.
	defer func() {
		if rollbackErr := pgxTx.Rollback(context.WithoutCancel(ctx)); rollbackErr != nil &&
			!errors.Is(rollbackErr, pgx.ErrTxClosed) {
			retErr = errors.Join(retErr, rollbackErr)
		}
	}()

	if err := fn(ctx, tx.Conn{Pgx: pgxTx}); err != nil {
		return err
	}

	return pgxTx.Commit(ctx)
}

func (t *transactorPgx) ImplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	if fn == nil {
		return errors.New("connection callback is nil")
	}

	conn, err := t.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return fn(ctx, tx.Conn{Pgx: conn})
}

// Compile-time assertion that pgxpool.Pool satisfies the narrow pool contract.
var _ pgxPool = (*pgxpool.Pool)(nil)
