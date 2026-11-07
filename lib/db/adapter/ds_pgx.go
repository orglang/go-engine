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

type transactorContextKey struct{}

func nested(ctx context.Context) bool {
	return ctx.Value(transactorContextKey{}) == true
}

func callbackContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(context.WithValue(ctx, transactorContextKey{}, true))
}

func (t *transactorPgx) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	if nested(ctx) {
		return tx.ErrNestedTransaction
	}
	if fn == nil {
		return errors.New("transaction callback is nil")
	}

	txCtx, cancel := callbackContext(ctx)
	defer cancel()

	pgxTx, err := t.pool.Begin(txCtx)
	if err != nil {
		return err
	}

	callbackErr := fn(txCtx, tx.Conn{Pgx: pgxTx})
	if callbackErr != nil {
		return errors.Join(callbackErr, pgxTx.Rollback(ctx))
	}
	if err := pgxTx.Commit(ctx); err != nil {
		return errors.Join(err, pgxTx.Rollback(ctx))
	}
	return nil
}

func (t *transactorPgx) ImplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	if nested(ctx) {
		return tx.ErrNestedTransaction
	}
	if fn == nil {
		return errors.New("connection callback is nil")
	}

	conn, err := t.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	txCtx, cancel := callbackContext(ctx)
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
var _ tx.Pgx = (pgx.Tx)(nil)
