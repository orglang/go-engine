package adapter

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"orglang/go-engine/lib/tx"
)

type stubPgx struct{}

func (stubPgx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (stubPgx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, nil
}

func (stubPgx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults {
	return nil
}

type stubTransaction struct {
	stubPgx
	commitErr   error
	rollbackErr error
	commitCalls int
	rollbackCalls int
	rollbackCtx context.Context
}

func (s *stubTransaction) Commit(context.Context) error {
	s.commitCalls++
	return s.commitErr
}

func (s *stubTransaction) Rollback(ctx context.Context) error {
	s.rollbackCalls++
	s.rollbackCtx = ctx
	return s.rollbackErr
}

type stubConnection struct {
	stubPgx
	releaseCalls int
}

func (s *stubConnection) Release() {
	s.releaseCalls++
}

type stubPool struct {
	transaction pgxTransaction
	connection  pgxConnection
	beginErr    error
	acquireErr  error
	beginCalls  int
	acquireCalls int
}

func (s *stubPool) Begin(context.Context) (pgxTransaction, error) {
	s.beginCalls++
	return s.transaction, s.beginErr
}

func (s *stubPool) Acquire(context.Context) (pgxConnection, error) {
	s.acquireCalls++
	return s.connection, s.acquireErr
}

func TestExplicitTxCommitsOnSuccess(t *testing.T) {
	txStub := &stubTransaction{rollbackErr: pgx.ErrTxClosed}
	pool := &stubPool{transaction: txStub}
	transactor := &transactorPgx{pool: pool}

	err := transactor.ExplicitTx(context.Background(), func(_ context.Context, conn tx.Conn) error {
		if conn.Pgx != txStub {
			t.Fatal("callback did not receive the transaction connection")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("ExplicitTx() error = %v, want nil", err)
	}
	if txStub.commitCalls != 1 {
		t.Errorf("commit calls = %d, want 1", txStub.commitCalls)
	}
	if txStub.rollbackCalls != 1 {
		t.Errorf("rollback calls = %d, want 1 (deferred safety rollback)", txStub.rollbackCalls)
	}
}

func TestExplicitTxRollsBackOnCallbackError(t *testing.T) {
	callbackErr := errors.New("callback failed")
	txStub := &stubTransaction{rollbackErr: pgx.ErrTxClosed}
	pool := &stubPool{transaction: txStub}
	transactor := &transactorPgx{pool: pool}

	err := transactor.ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
		return callbackErr
	})

	if !errors.Is(err, callbackErr) {
		t.Fatalf("ExplicitTx() error = %v, want callback error", err)
	}
	if txStub.commitCalls != 0 {
		t.Errorf("commit calls = %d, want 0", txStub.commitCalls)
	}
	if txStub.rollbackCalls != 1 {
		t.Errorf("rollback calls = %d, want 1", txStub.rollbackCalls)
	}
}

func TestExplicitTxReturnsCommitError(t *testing.T) {
	commitErr := errors.New("commit failed")
	txStub := &stubTransaction{commitErr: commitErr, rollbackErr: pgx.ErrTxClosed}
	pool := &stubPool{transaction: txStub}
	transactor := &transactorPgx{pool: pool}

	err := transactor.ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
		return nil
	})

	if !errors.Is(err, commitErr) {
		t.Fatalf("ExplicitTx() error = %v, want commit error", err)
	}
	if txStub.commitCalls != 1 {
		t.Errorf("commit calls = %d, want 1", txStub.commitCalls)
	}
	if txStub.rollbackCalls != 1 {
		t.Errorf("rollback calls = %d, want 1", txStub.rollbackCalls)
	}
}

func TestExplicitTxJoinsCallbackAndRollbackErrors(t *testing.T) {
	callbackErr := errors.New("callback failed")
	rollbackErr := errors.New("rollback failed")
	txStub := &stubTransaction{rollbackErr: rollbackErr}
	pool := &stubPool{transaction: txStub}
	transactor := &transactorPgx{pool: pool}

	err := transactor.ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
		return callbackErr
	})

	if !errors.Is(err, callbackErr) {
		t.Errorf("ExplicitTx() error = %v, want callback error", err)
	}
	if !errors.Is(err, rollbackErr) {
		t.Errorf("ExplicitTx() error = %v, want rollback error", err)
	}
	if txStub.rollbackCalls != 1 {
		t.Errorf("rollback calls = %d, want 1", txStub.rollbackCalls)
	}
}

func TestExplicitTxReturnsRollbackErrorWithoutCallbackError(t *testing.T) {
	rollbackErr := errors.New("rollback failed")
	txStub := &stubTransaction{rollbackErr: rollbackErr}
	pool := &stubPool{transaction: txStub}
	transactor := &transactorPgx{pool: pool}

	err := transactor.ExplicitTx(context.Background(), func(context.Context, tx.Conn) error {
		return nil
	})

	if !errors.Is(err, rollbackErr) {
		t.Fatalf("ExplicitTx() error = %v, want rollback error", err)
	}
}

func TestExplicitTxRollbackIgnoresCanceledParentContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	txStub := &stubTransaction{rollbackErr: pgx.ErrTxClosed}
	pool := &stubPool{transaction: txStub}
	transactor := &transactorPgx{pool: pool}

	_ = transactor.ExplicitTx(ctx, func(callbackCtx context.Context, _ tx.Conn) error {
		if callbackCtx.Err() == nil {
			t.Error("callback context error = nil, want canceled context")
		}
		return errors.New("callback failed")
	})

	if txStub.rollbackCtx == nil {
		t.Fatal("rollback context was nil")
	}
	if txStub.rollbackCtx.Err() != nil {
		t.Errorf("rollback context error = %v, want nil", txStub.rollbackCtx.Err())
	}
}

func TestImplicitTxUsesAndReleasesOnePooledConnection(t *testing.T) {
	conn := &stubConnection{}
	pool := &stubPool{connection: conn}
	transactor := &transactorPgx{pool: pool}
	callbackErr := errors.New("callback failed")

	err := transactor.ImplicitTx(context.Background(), func(_ context.Context, got tx.Conn) error {
		if got.Pgx != conn {
			t.Fatal("callback did not receive the acquired connection")
		}
		if conn.releaseCalls != 0 {
			t.Error("connection was released before callback returned")
		}
		return callbackErr
	})

	if !errors.Is(err, callbackErr) {
		t.Fatalf("ImplicitTx() error = %v, want callback error", err)
	}
	if pool.acquireCalls != 1 {
		t.Errorf("acquire calls = %d, want 1", pool.acquireCalls)
	}
	if pool.beginCalls != 0 {
		t.Errorf("begin calls = %d, want 0 (ImplicitTx must not start a transaction)", pool.beginCalls)
	}
	if conn.releaseCalls != 1 {
		t.Errorf("release calls = %d, want 1", conn.releaseCalls)
	}
}
