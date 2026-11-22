package adapter

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/tx"
)

type batchQueryErrorPgx struct {
	batch pgx.BatchResults
}

func (p batchQueryErrorPgx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	panic("unexpected Exec")
}

func (p batchQueryErrorPgx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unexpected Query")
}

func (p batchQueryErrorPgx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults {
	return p.batch
}

type batchQueryErrorResults struct {
	queryErr error
	closed   bool
}

func (r *batchQueryErrorResults) Exec() (pgconn.CommandTag, error) {
	panic("unexpected batch Exec")
}

func (r *batchQueryErrorResults) Query() (pgx.Rows, error) {
	return nil, r.queryErr
}

func (r *batchQueryErrorResults) QueryRow() pgx.Row {
	panic("unexpected batch QueryRow")
}

func (r *batchQueryErrorResults) Close() error {
	r.closed = true
	return nil
}

func (r *batchQueryErrorResults) Tx() pgx.Tx {
	panic("unexpected batch Tx")
}

func TestGetRecsReturnsBatchQueryError(t *testing.T) {
	queryErr := errors.New("batch query failed")
	batch := &batchQueryErrorResults{queryErr: queryErr}
	dao := &daoEnt{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	conn := tx.Conn{Pgx: batchQueryErrorPgx{batch: batch}}

	_, err := dao.GetRecs(context.Background(), conn, []identity.ADT{identity.New()})
	if !errors.Is(err, queryErr) {
		t.Fatalf("GetRecs() error = %v, want %v", err, queryErr)
	}
	if !batch.closed {
		t.Fatal("GetRecs() did not close batch results")
	}
}
