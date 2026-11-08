package adapter

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/tx"
)

type testPgxPool struct {
	tx        pgx.Tx
	beginErr  error
	beginCalls int
}

func (p *testPgxPool) Begin(context.Context) (pgx.Tx, error) {
	p.beginCalls++
	return p.tx, p.beginErr
}

func (*testPgxPool) Acquire(context.Context) (*pgxpool.Conn, error) {
	panic("Acquire should not be called by ExplicitTx tests")
}

type testPgxTx struct {
	pgx.Tx
	commitErr  error
	rollbackErr error
	commitCalls int
	rollbackCalls int
	rollbackContextErr error
}

func (t *testPgxTx) Commit(context.Context) error {
	t.commitCalls++
	return t.commitErr
}

func (t *testPgxTx) Rollback(ctx context.Context) error {
	t.rollbackCalls++
	t.rollbackContextErr = ctx.Err()
	return t.rollbackErr
}

func TestTransactorPgxExplicitTx(t *testing.T) {
	t.Run("nil callback", func(t *testing.T) {
		pool := &testPgxPool{}
		err := (&transactorPgx{pool: pool}).ExplicitTx(context.Background(), nil)
		if err == nil {
			t.Fatal("expected an error for a nil callback")
		}
		if pool.beginCalls != 0 {
			t.Fatalf("Begin called %d times, want 0", pool.beginCalls)
		}
	})

	t.Run("begin error", func(t *testing.T) {
		beginErr := errors.New("begin failed")
		pool := &testPgxPool{beginErr: beginErr}
		err := (&transactorPgx{pool: pool}).ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
			t.Fatal("callback must not be called")
			return nil
		})
		if !errors.Is(err, beginErr) {
			t.Fatalf("got error %v, want %v", err, beginErr)
		}
		if pool.beginCalls != 1 {
			t.Fatalf("Begin called %d times, want 1", pool.beginCalls)
		}
	})

	t.Run("callback error rolls back", func(t *testing.T) {
		callbackErr := errors.New("callback failed")
		pgxTx := &testPgxTx{}
		pool := &testPgxPool{tx: pgxTx}
		err := (&transactorPgx{pool: pool}).ExplicitTx(context.Background(), func(_ context.Context, conn tx.Conn) error {
			if conn.Pgx != pgxTx {
				t.Fatal("callback received a different transaction")
			}
			return callbackErr
		})
		if !errors.Is(err, callbackErr) {
			t.Fatalf("got error %v, want %v", err, callbackErr)
		}
		if pgxTx.commitCalls != 0 {
			t.Fatalf("Commit called %d times, want 0", pgxTx.commitCalls)
		}
		if pgxTx.rollbackCalls != 1 {
			t.Fatalf("Rollback called %d times, want 1", pgxTx.rollbackCalls)
		}
	})

	t.Run("callback and rollback errors are joined", func(t *testing.T) {
		callbackErr := errors.New("callback failed")
		rollbackErr := errors.New("rollback failed")
		pgxTx := &testPgxTx{rollbackErr: rollbackErr}
		pool := &testPgxPool{tx: pgxTx}
		err := (&transactorPgx{pool: pool}).ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
			return callbackErr
		})
		if !errors.Is(err, callbackErr) || !errors.Is(err, rollbackErr) {
			t.Fatalf("got error %v, want callback and rollback errors", err)
		}
		if pgxTx.rollbackCalls != 1 {
			t.Fatalf("Rollback called %d times, want 1", pgxTx.rollbackCalls)
		}
	})

	t.Run("success commits and ignores closed transaction rollback", func(t *testing.T) {
		pgxTx := &testPgxTx{rollbackErr: pgx.ErrTxClosed}
		pool := &testPgxPool{tx: pgxTx}
		callbackCtx := context.WithValue(context.Background(), struct{}{}, "callback")
		err := (&transactorPgx{pool: pool}).ExplicitTx(callbackCtx, func(ctx context.Context, conn tx.Conn) error {
			if ctx != callbackCtx {
				t.Fatal("callback did not receive the original context")
			}
			if conn.Pgx != pgxTx {
				t.Fatal("callback received a different transaction")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("ExplicitTx returned unexpected error: %v", err)
		}
		if pgxTx.commitCalls != 1 {
			t.Fatalf("Commit called %d times, want 1", pgxTx.commitCalls)
		}
		if pgxTx.rollbackCalls != 1 {
			t.Fatalf("Rollback called %d times, want 1", pgxTx.rollbackCalls)
		}
	})

	t.Run("commit error is returned", func(t *testing.T) {
		commitErr := errors.New("commit failed")
		pgxTx := &testPgxTx{commitErr: commitErr, rollbackErr: pgx.ErrTxClosed}
		pool := &testPgxPool{tx: pgxTx}
		err := (&transactorPgx{pool: pool}).ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
			return nil
		})
		if !errors.Is(err, commitErr) {
			t.Fatalf("got error %v, want %v", err, commitErr)
		}
		if pgxTx.rollbackCalls != 1 {
			t.Fatalf("Rollback called %d times, want 1", pgxTx.rollbackCalls)
		}
	})

	t.Run("rollback uses context without parent cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		pgxTx := &testPgxTx{rollbackErr: pgx.ErrTxClosed}
		pool := &testPgxPool{tx: pgxTx}
		err := (&transactorPgx{pool: pool}).ExplicitTx(ctx, func(context.Context, tx.Conn) error {
			return errors.New("callback failed")
		})
		if err == nil {
			t.Fatal("expected callback error")
		}
		if pgxTx.rollbackContextErr != nil {
			t.Fatalf("rollback context is canceled: %v", pgxTx.rollbackContextErr)
		}
	})
}
