package adapter

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"orglang/go-engine/internal/adt/uniqsym"
	dbadapter "orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/pool/ent"
	progrootadapter "orglang/go-engine/internal/pool/progroot/adapter"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	"orglang/go-engine/internal/testutil/pgxtest"
)

func TestIntegrationTypedefCreateRollsBackAllRepositories(t *testing.T) {
	pool := pgxtest.NewPostgresPool(t)

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection: %v", err)
	}
	_, err = conn.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS ltree;
		CREATE TEMP TABLE pool_prog_roots (prog_id varchar PRIMARY KEY, prog_rn bigint NOT NULL);
		CREATE TEMP TABLE pool_prog_binds (
			prog_qn ltree UNIQUE,
			prog_id varchar,
			kind smallint
		);
		CREATE TEMP TABLE pool_type_exprs (
			expr_ck bigint UNIQUE,
			sup_expr_ck bigint,
			kind smallint,
			spec jsonb
		);
		CREATE TEMP TABLE pool_type_def_guard (id varchar PRIMARY KEY);
		CREATE TEMP TABLE pool_type_defs (
			type_id varchar UNIQUE REFERENCES pool_type_def_guard(id)
				DEFERRABLE INITIALLY DEFERRED,
			prog_id varchar NOT NULL REFERENCES pool_prog_roots(prog_id),
			expr_ck bigint
		);
	`)
	if err != nil {
		conn.Release()
		t.Fatalf("create temporary integration schema: %v", err)
	}
	conn.Release()

	logger := slog.Default()
	entClient := ent.NewClient(ent.Driver(entpgx.New()))
	t.Cleanup(func() {
		if err := entClient.Close(); err != nil {
			t.Errorf("close Ent client: %v", err)
		}
	})

	typeDefRepo := NewDaoEnt(entClient, logger)
	service := typedef.NewService(
		typeDefRepo,
		progrootadapter.NewDaoEnt(entClient, logger),
		dbadapter.NewTransactorPgx(pool),
		logger,
	)

	typeQN, err := uniqsym.ConvertFromString("integration.atomicity")
	if err != nil {
		t.Fatalf("create type qualified name: %v", err)
	}
	_, _, err = service.Create(ctx, typedef.DefSpec{
		TypeQN:  typeQN,
		TypeExprr: typedef.OneSpec{},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want deferred constraint failure at commit")
	}

	verifyConn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection for verification: %v", err)
	}
	defer verifyConn.Release()

	for _, table := range []string{"pool_prog_roots", "pool_prog_binds", "pool_type_exprs", "pool_type_defs"} {
		var count int
		if err := verifyConn.QueryRow(ctx, "SELECT count(*) FROM "+table).Scan(&count); err != nil {
			t.Fatalf("count rows in %s: %v", table, err)
		}
		if count != 0 {
			t.Errorf("%s contains %d rows after failed Create(), want 0", table, count)
		}
	}
}

func TestIntegrationTypedefCreatePersistsAllRepositories(t *testing.T) {
	pool := pgxtest.NewPostgresPool(t)

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection: %v", err)
	}
	_, err = conn.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS ltree;
		CREATE TEMP TABLE pool_prog_roots (prog_id varchar PRIMARY KEY, prog_rn bigint NOT NULL);
		CREATE TEMP TABLE pool_prog_binds (
			prog_qn ltree UNIQUE,
			prog_id varchar,
			kind smallint
		);
		CREATE TEMP TABLE pool_type_exprs (
			expr_ck bigint UNIQUE,
			sup_expr_ck bigint,
			kind smallint,
			spec jsonb
		);
		CREATE TEMP TABLE pool_type_defs (
			type_id varchar UNIQUE,
			prog_id varchar NOT NULL REFERENCES pool_prog_roots(prog_id),
			expr_ck bigint
		);
	`)
	if err != nil {
		conn.Release()
		t.Fatalf("create temporary integration schema: %v", err)
	}
	conn.Release()

	logger := slog.Default()
	entClient := ent.NewClient(ent.Driver(entpgx.New()))
	t.Cleanup(func() {
		if err := entClient.Close(); err != nil {
			t.Errorf("close Ent client: %v", err)
		}
	})

	service := typedef.NewService(
		NewDaoEnt(entClient, logger),
		progrootadapter.NewDaoEnt(entClient, logger),
		dbadapter.NewTransactorPgx(pool),
		logger,
	)

	typeQN, err := uniqsym.ConvertFromString("integration.atomicity.success")
	if err != nil {
		t.Fatalf("create type qualified name: %v", err)
	}
	typeID, progRef, err := service.Create(ctx, typedef.DefSpec{
		TypeQN:  typeQN,
		TypeExprr: typedef.OneSpec{},
	})
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}
	if typeID.IsEmpty() {
		t.Fatal("Create() returned an empty type ID")
	}

	verifyConn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection for verification: %v", err)
	}
	defer verifyConn.Release()

	var descBindCount, progRootCount, typeExpCount, typeDefCount int
	err = verifyConn.QueryRow(ctx,
		"SELECT count(*) FROM pool_prog_binds WHERE prog_id = $1", typeID.String(),
	).Scan(&descBindCount)
	if err != nil {
		t.Fatalf("count descriptor bindings: %v", err)
	}
	err = verifyConn.QueryRow(ctx, "SELECT count(*) FROM pool_prog_roots WHERE prog_id = $1", progRef.ProgID.String()).Scan(&progRootCount)
	if err != nil {
		t.Fatalf("count Program roots: %v", err)
	}
	err = verifyConn.QueryRow(ctx, "SELECT count(*) FROM pool_type_exprs").Scan(&typeExpCount)
	if err != nil {
		t.Fatalf("count type expressions: %v", err)
	}
	err = verifyConn.QueryRow(ctx,
		"SELECT count(*) FROM pool_type_defs WHERE type_id = $1 AND prog_id = $2", typeID.String(), progRef.ProgID.String(),
	).Scan(&typeDefCount)
	if err != nil {
		t.Fatalf("count type definitions: %v", err)
	}

	if descBindCount != 1 || progRootCount != 1 || typeExpCount == 0 || typeDefCount != 1 {
		t.Errorf(
			"persisted rows: desc binds = %d, desc roots = %d, type exps = %d, type defs = %d; want 1, 1, >0, 1",
			descBindCount, progRootCount, typeExpCount, typeDefCount,
		)
	}
}
