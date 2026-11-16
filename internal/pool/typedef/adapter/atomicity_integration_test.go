package adapter

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	dbadapter "orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/pool/ent"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	typeexpadapter "orglang/go-engine/internal/pool/typeexp/adapter"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

type integrationTypeExpQueryBuilder struct{}

func (integrationTypeExpQueryBuilder) InsertRec(st typeexp.StateDS) (string, []any) {
	return "INSERT INTO pool_type_exps (exp_ck, sup_exp_ck, kind, spec) VALUES ($1, $2, $3, $4)",
		[]any{st.ExpCK, st.SupExpCK, st.K, []byte("{}")}
}

func (integrationTypeExpQueryBuilder) SelectRecByCK() string {
	return ""
}

func TestIntegrationTypedefCreateRollsBackAllRepositories(t *testing.T) {
	dsn := os.Getenv("ENTPGX_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("ENTPGX_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse PostgreSQL config: %v", err)
	}
	config.MaxConns = 1
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("open PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection: %v", err)
	}
	_, err = conn.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS ltree;
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
			type_rn bigint,
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
	typeExpRepo := typeexpadapter.NewDaoPgx(integrationTypeExpQueryBuilder{}, logger)
	descSemRepo := descsemadapter.NewDaoPgx(descBinds)(logger)
	service := typedef.NewService(
		typeDefRepo,
		typeExpRepo,
		descSemRepo,
		dbadapter.NewTransactorPgx(pool),
		logger,
	)

	typeQN, err := uniqsym.ConvertFromString("integration.atomicity")
	if err != nil {
		t.Fatalf("create type qualified name: %v", err)
	}
	_, err = service.Create(typedef.DefSpec{
		TypeQN:  typeQN,
		TypeExp: typeexp.OneSpec{},
	})
	if err == nil {
		t.Fatal("Create() error = nil, want deferred constraint failure at commit")
	}

	verifyConn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection for verification: %v", err)
	}
	defer verifyConn.Release()

	for _, table := range []string{"pool_desc_binds", "pool_type_exps", "pool_type_defs"} {
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
	dsn := os.Getenv("ENTPGX_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("ENTPGX_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse PostgreSQL config: %v", err)
	}
	config.MaxConns = 1
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("open PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection: %v", err)
	}
	_, err = conn.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS ltree;
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
			type_rn bigint,
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
		typeexpadapter.NewDaoPgx(integrationTypeExpQueryBuilder{}, logger),
		descsemadapter.NewDaoPgx(descBinds)(logger),
		dbadapter.NewTransactorPgx(pool),
		logger,
	)

	typeQN, err := uniqsym.ConvertFromString("integration.atomicity.success")
	if err != nil {
		t.Fatalf("create type qualified name: %v", err)
	}
	ref, err := service.Create(typedef.DefSpec{
		TypeQN:  typeQN,
		TypeExp: typeexp.OneSpec{},
	})
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}
	if ref.TypeID == "" {
		t.Fatal("Create() returned an empty type ID")
	}

	verifyConn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection for verification: %v", err)
	}
	defer verifyConn.Release()

	var descBindCount, typeExpCount, typeDefCount int
	if err := verifyConn.QueryRow(ctx,
		"SELECT count(*) FROM pool_desc_binds WHERE desc_id = $1", ref.TypeID,
	).Scan(&descBindCount); err != nil {
		t.Fatalf("count descriptor bindings: %v", err)
	}
	if err := verifyConn.QueryRow(ctx, "SELECT count(*) FROM pool_type_exps").
		Scan(&typeExpCount); err != nil {
		t.Fatalf("count type expressions: %v", err)
	}
	if err := verifyConn.QueryRow(ctx,
		"SELECT count(*) FROM pool_type_defs WHERE type_id = $1", ref.TypeID,
	).Scan(&typeDefCount); err != nil {
		t.Fatalf("count type definitions: %v", err)
	}

	if descBindCount != 1 || typeExpCount == 0 || typeDefCount != 1 {
		t.Errorf(
			"persisted rows: desc binds = %d, type exps = %d, type defs = %d; want 1, >0, 1",
			descBindCount, typeExpCount, typeDefCount,
		)
	}
}
