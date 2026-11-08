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

func (t *transactorPgx) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	if fn == nil {
		return errors.New("transaction callback is nil")
	}

	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	pgxTx, err := t.pool.Begin(txCtx)
	if err != nil {
		return err
	}

	// Always attempt rollback on exit. pgx returns ErrTxClosed after a
	// successful commit or an earlier rollback; that cleanup error is ignored.
	defer func() {
		_ = pgxTx.Rollback(context.WithoutCancel(ctx))
	}()

	callbackErr := fn(txCtx, tx.Conn{Pgx: pgxTx})
	if callbackErr != nil {
		return errors.Join(callbackErr, pgxTx.Rollback(context.WithoutCancel(ctx)))
	}
	if err := pgxTx.Commit(ctx); err != nil {
		return errors.Join(err, pgxTx.Rollback(context.WithoutCancel(ctx)))
	}
	return nil
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

	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fn(txCtx, tx.Conn{Pgx: conn})
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
