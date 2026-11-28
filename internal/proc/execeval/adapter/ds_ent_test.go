package adapter

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/lib/tx"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

type batchExecErrorPgx struct {
	batch pgx.BatchResults
}

func (p batchExecErrorPgx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	panic("unexpected Eval")
}

func (p batchExecErrorPgx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unexpected Query")
}

func (p batchExecErrorPgx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults {
	return p.batch
}

type batchExecErrorResults struct {
	execErr error
	closed  bool
}

func (r *batchExecErrorResults) Exec() (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, r.execErr
}

func (r *batchExecErrorResults) Query() (pgx.Rows, error) {
	panic("unexpected batch Query")
}

func (r *batchExecErrorResults) QueryRow() pgx.Row {
	panic("unexpected batch QueryRow")
}

func (r *batchExecErrorResults) Close() error {
	r.closed = true
	return nil
}

func (r *batchExecErrorResults) Tx() pgx.Tx {
	panic("unexpected batch Tx")
}

func TestModifyRecReturnsBatchExecError(t *testing.T) {
	execErr := errors.New("batch execution failed")
	batch := &batchExecErrorResults{execErr: execErr}
	dao := &daoEnt{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	conn := tx.Conn{Pgx: batchExecErrorPgx{batch: batch}}

	err := dao.ModifyRec(context.Background(), conn, execeval.EvalMod{
		RootRefs: []execroot.RootRef{execroot.New()},
	})
	if !errors.Is(err, execErr) {
		t.Fatalf("ModifyRec() error = %v, want %v", err, execErr)
	}
	if !batch.closed {
		t.Fatal("ModifyRec() did not close batch results")
	}
}
