// Package adapter provides the pgx implementation of the transaction contract.
package adapter

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"orglang/go-engine/lib/tx"
)

type pgxTransaction interface {
	tx.Pgx
	Commit(context.Context) error
	Rollback(context.Context) error
}

type pgxConnection interface {
	tx.Pgx
	Release()
}

type pgxPool interface {
	Begin(context.Context) (pgxTransaction, error)
	Acquire(context.Context) (pgxConnection, error)
}

type pgxPoolAdapter struct {
	pool *pgxpool.Pool
}

func (p pgxPoolAdapter) Begin(ctx context.Context) (pgxTransaction, error) {
	return p.pool.Begin(ctx)
}

func (p pgxPoolAdapter) Acquire(ctx context.Context) (pgxConnection, error) {
	return p.pool.Acquire(ctx)
}

type transactorPgx struct {
	pool pgxPool
}

// NewTransactorPgx creates a PostgreSQL-backed transactor.
func NewTransactorPgx(pool *pgxpool.Pool) tx.Transactor {
	return &transactorPgx{pool: pgxPoolAdapter{pool: pool}}
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

// Compile-time assertions for the concrete pgx resources and command-tag type.
var (
	_ tx.Pgx           = (pgx.Tx)(nil)
	_ tx.Pgx           = (*pgxpool.Conn)(nil)
	_ pgxTransaction   = (pgx.Tx)(nil)
	_ pgxConnection    = (*pgxpool.Conn)(nil)
	_ pgconn.CommandTag = pgconn.CommandTag{}
)
