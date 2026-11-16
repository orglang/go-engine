package entpgx

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"orglang/go-engine/internal/pool/ent"
)

func TestPostgresIntegrationUsesCallbackTransaction(t *testing.T) {
	dsn := os.Getenv("ENTPGX_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("ENTPGX_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("open PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	for _, test := range []struct {
		name    string
		commit  bool
		wantRow int
	}{
		{name: "rollback", commit: false, wantRow: 0},
		{name: "commit", commit: true, wantRow: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			conn, err := pool.Acquire(ctx)
			if err != nil {
				t.Fatalf("acquire PostgreSQL connection: %v", err)
			}
			t.Cleanup(conn.Release)

			_, err = conn.Exec(ctx, "CREATE TEMP TABLE pool_type_defs (type_id varchar UNIQUE, type_rn bigint, exp_ck bigint) ON COMMIT PRESERVE ROWS")
			if err != nil {
				t.Fatalf("create temporary table: %v", err)
			}
			t.Cleanup(func() {
				_, dropErr := conn.Exec(context.Background(), "DROP TABLE IF EXISTS pool_type_defs")
				if dropErr != nil {
					t.Errorf("drop temporary table: %v", dropErr)
				}
			})

			tx, err := conn.Begin(ctx)
			if err != nil {
				t.Fatalf("begin transaction: %v", err)
			}
			defer func() {
				_ = tx.Rollback(context.Background())
			}()

			client := ent.NewClient(ent.Driver(New()))
			t.Cleanup(func() {
				if err := client.Close(); err != nil {
					t.Errorf("close Ent client: %v", err)
				}
			})

			id := "entpgx-" + test.name
			entCtx := WithConn(ctx, tx)
			created, err := client.TypeDef.Create().
				SetID(id).
				SetTypeRn(7).
				SetExpCk(11).
				Save(entCtx)
			if err != nil {
				t.Fatalf("create entity through Ent: %v", err)
			}
			if created.ID != id {
				t.Fatalf("created ID = %q, want %q", created.ID, id)
			}

			loaded, err := client.TypeDef.Get(entCtx, id)
			if err != nil {
				t.Fatalf("read entity through Ent in callback transaction: %v", err)
			}
			if loaded.TypeRn != 7 || loaded.ExpCk != 11 {
				t.Fatalf("loaded entity = %#v, want type_rn=7 exp_ck=11", loaded)
			}

			var visibleID string
			err = tx.QueryRow(ctx, "SELECT type_id FROM pool_type_defs WHERE type_id = $1", id).Scan(&visibleID)
			if err != nil {
				t.Fatalf("read Ent insert through native pgx transaction: %v", err)
			}
			if visibleID != id {
				t.Fatalf("pgx read ID = %q, want %q", visibleID, id)
			}

			if test.commit {
				if err := tx.Commit(ctx); err != nil {
					t.Fatalf("commit transaction: %v", err)
				}
			} else {
				if err := tx.Rollback(ctx); err != nil {
					t.Fatalf("rollback transaction: %v", err)
				}
			}

			var count int
			err = conn.QueryRow(ctx, "SELECT count(*) FROM pool_type_defs WHERE type_id = $1", id).Scan(&count)
			if err != nil {
				t.Fatalf("count rows after transaction: %v", err)
			}
			if count != test.wantRow {
				t.Fatalf("row count after %s = %d, want %d", test.name, count, test.wantRow)
			}
		})
	}
}
