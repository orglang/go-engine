package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"orglang/go-engine/internal/adt/contkey"
	dbadapter "orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	typedef "orglang/go-engine/internal/proc/typedef/core"
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

	repo := typedefadapter.NewExpDaoEnt(client, slog.Default())
	transactor := dbadapter.NewTransactorPgx(pool)
	want := typedef.TensorRec{
		ExpCK: contkey.Three,
		Val:   typedef.OneRec{ExpCK: contkey.One},
		Cont:  typedef.OneRec{ExpCK: contkey.Two},
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
		rows, err := conn.Pgx.Query(txCtx, "SELECT count(*) FROM proc_type_exps")
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err
			}
			return errors.New("count query returned no rows")
		}
		var count int
		if err := rows.Scan(&count); err != nil {
			return err
		}
		if count != 3 {
			t.Errorf("stored state count = %d, want 3 after duplicate insert", count)
		}
		rows.Close()
		return rows.Err()
	})
	if err != nil {
		t.Fatalf("verify persisted type expression: %v", err)
	}
}
