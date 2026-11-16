package adapter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	dbadapter "orglang/go-engine/internal/lib/db/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/pool/ent"
	typeexpadapter "orglang/go-engine/internal/pool/typeexp/adapter"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
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
		CREATE TEMP TABLE pool_type_defs (
			type_id varchar UNIQUE,
			type_rn bigint,
			exp_ck bigint
		);
		CREATE FUNCTION pg_temp.reject_type_def_commit() RETURNS trigger
		LANGUAGE plpgsql AS $$
		BEGIN
			RAISE EXCEPTION 'forced deferred constraint failure';
		END;
		$$;
		CREATE CONSTRAINT TRIGGER reject_type_def_commit
		AFTER INSERT ON pool_type_defs
		DEFERRABLE INITIALLY DEFERRED
		FOR EACH ROW EXECUTE FUNCTION pg_temp.reject_type_def_commit();
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
