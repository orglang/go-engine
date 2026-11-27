package adapter

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"orglang/go-engine/internal/adt/uniqsym"
	dbadapter "orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	descrootadapter "orglang/go-engine/internal/pool/descroot/adapter"
	"orglang/go-engine/internal/pool/ent"
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
		CREATE TEMP TABLE pool_desc_roots (desc_id varchar PRIMARY KEY, desc_rn bigint NOT NULL);
		CREATE TEMP TABLE pool_desc_binds (
			desc_qn ltree UNIQUE,
			desc_id varchar,
			kind smallint
		);
		CREATE TEMP TABLE pool_type_exps (
			exp_ck bigint UNIQUE,
			sup_exp_ck bigint,
			kind smallint,
			spec jsonb
		);
		CREATE TEMP TABLE pool_type_def_guard (id varchar PRIMARY KEY);
		CREATE TEMP TABLE pool_type_defs (
			type_id varchar UNIQUE REFERENCES pool_type_def_guard(id)
				DEFERRABLE INITIALLY DEFERRED,
			desc_id varchar NOT NULL REFERENCES pool_desc_roots(desc_id),
			exp_ck bigint
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
		descrootadapter.NewDaoEnt(entClient, logger),
		dbadapter.NewTransactorPgx(pool),
		logger,
	)

	typeQN, err := uniqsym.ConvertFromString("integration.atomicity")
	if err != nil {
		t.Fatalf("create type qualified name: %v", err)
	}
	_, err = service.Create(ctx, typedef.DefSpec{
		TypeQN:  typeQN,
		TypeExp: typedef.OneSpec{},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want deferred constraint failure at commit")
	}

	verifyConn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection for verification: %v", err)
	}
	defer verifyConn.Release()

	for _, table := range []string{"pool_desc_roots", "pool_desc_binds", "pool_type_exps", "pool_type_defs"} {
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
		CREATE TEMP TABLE pool_desc_roots (desc_id varchar PRIMARY KEY, desc_rn bigint NOT NULL);
		CREATE TEMP TABLE pool_desc_binds (
			desc_qn ltree UNIQUE,
			desc_id varchar,
			kind smallint
		);
		CREATE TEMP TABLE pool_type_exps (
			exp_ck bigint UNIQUE,
			sup_exp_ck bigint,
			kind smallint,
			spec jsonb
		);
		CREATE TEMP TABLE pool_type_defs (
			type_id varchar UNIQUE,
			desc_id varchar NOT NULL REFERENCES pool_desc_roots(desc_id),
			exp_ck bigint
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
		descrootadapter.NewDaoEnt(entClient, logger),
		dbadapter.NewTransactorPgx(pool),
		logger,
	)

	typeQN, err := uniqsym.ConvertFromString("integration.atomicity.success")
	if err != nil {
		t.Fatalf("create type qualified name: %v", err)
	}
	typeID, descRef, err := service.Create(ctx, typedef.DefSpec{
		TypeQN:  typeQN,
		TypeExp: typedef.OneSpec{},
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

	var descBindCount, descRootCount, typeExpCount, typeDefCount int
	err = verifyConn.QueryRow(ctx,
		"SELECT count(*) FROM pool_desc_binds WHERE desc_id = $1", typeID.String(),
	).Scan(&descBindCount)
	if err != nil {
		t.Fatalf("count descriptor bindings: %v", err)
	}
	err = verifyConn.QueryRow(ctx, "SELECT count(*) FROM pool_desc_roots WHERE desc_id = $1", descRef.DescID.String()).Scan(&descRootCount)
	if err != nil {
		t.Fatalf("count Description roots: %v", err)
	}
	err = verifyConn.QueryRow(ctx, "SELECT count(*) FROM pool_type_exps").Scan(&typeExpCount)
	if err != nil {
		t.Fatalf("count type expressions: %v", err)
	}
	err = verifyConn.QueryRow(ctx,
		"SELECT count(*) FROM pool_type_defs WHERE type_id = $1 AND desc_id = $2", typeID.String(), descRef.DescID.String(),
	).Scan(&typeDefCount)
	if err != nil {
		t.Fatalf("count type definitions: %v", err)
	}

	if descBindCount != 1 || descRootCount != 1 || typeExpCount == 0 || typeDefCount != 1 {
		t.Errorf(
			"persisted rows: desc binds = %d, desc roots = %d, type exps = %d, type defs = %d; want 1, 1, >0, 1",
			descBindCount, descRootCount, typeExpCount, typeDefCount,
		)
	}
}
