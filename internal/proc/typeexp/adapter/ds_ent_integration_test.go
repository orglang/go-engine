package adapter

import (
	"context"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"orglang/go-engine/internal/adt/contkey"
	dbadapter "orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
	"orglang/go-engine/internal/testutil/pgxtest"
)

func TestIntegrationProcTypeExpAddRecIsIdempotent(t *testing.T) {
	pool := pgxtest.NewPostgresPool(t)

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection: %v", err)
	}
	_, err = conn.Exec(ctx, `
		CREATE TEMP TABLE proc_type_exps (
			exp_ck bigint UNIQUE,
			sup_exp_ck bigint,
			kind smallint,
			spec jsonb
		);
	`)
	if err != nil {
		conn.Release()
		t.Fatalf("create temporary integration schema: %v", err)
	}
	conn.Release()

	client := ent.NewClient(ent.Driver(entpgx.New()))
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("close Ent client: %v", err)
		}
	})

	repo := NewDaoEnt(client, slog.Default())
	transactor := dbadapter.NewTransactorPgx(pool)
	want := typeexp.TensorRec{
		ExpCK: contkey.Three,
		Val:   typeexp.OneRec{ExpCK: contkey.One},
		Cont:  typeexp.OneRec{ExpCK: contkey.Two},
	}

	for range 2 {
		err := transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			return repo.AddRec(txCtx, conn, want)
		})
		if err != nil {
			t.Fatalf("AddRec() error = %v, want nil on repeated insert", err)
		}
	}

	err = transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		got, err := repo.SelectRecByCK(txCtx, conn, contkey.Three)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("SelectRecByCK() = %#v, want %#v", got, want)
		}
		var count int
		if err := conn.Pgx.QueryRow(txCtx, "SELECT count(*) FROM proc_type_exps").Scan(&count); err != nil {
			return err
		}
		if count != 3 {
			t.Errorf("stored state count = %d, want 3 after duplicate insert", count)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("verify persisted type expression: %v", err)
	}
}
