// Package adapter provides the pgx implementation of the transaction contract.
package adapter

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"orglang/go-engine/lib/tx"
)

type transactorPgx struct {
	pool *pgxpool.Pool
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

func (t *transactorPgx) ExplicitTxOld(ctx context.Context, fn func(tx.UoW) error) error {
	if fn == nil {
		return errors.New("transaction callback is nil")
	}
	return t.ExplicitTx(ctx, func(callbackCtx context.Context, conn tx.Conn) error {
		return fn(tx.UoW{Ctx: callbackCtx, Pgx: conn.Pgx})
	})
}

func (t *transactorPgx) ImplicitTxOld(ctx context.Context, fn func(tx.UoW) error) error {
	if fn == nil {
		return errors.New("connection callback is nil")
	}
	return t.ImplicitTx(ctx, func(callbackCtx context.Context, conn tx.Conn) error {
		return fn(tx.UoW{Ctx: callbackCtx, Pgx: conn.Pgx})
	})
}

// Compile-time assertion that the concrete driver transaction satisfies the
// minimal interface used by Conn.
var (
	_ tx.Pgx = (pgx.Tx)(nil)
	_ tx.Pgx = (*pgxpool.Conn)(nil)
)
