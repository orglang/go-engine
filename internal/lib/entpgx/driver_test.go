package entpgx

import (
	"context"
	"errors"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"orglang/go-engine/internal/lib/tx"
)

type fakeConn struct {
	tx.Pgx
	tag  pgconn.CommandTag
	rows pgx.Rows
}

func (c *fakeConn) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return c.tag, nil
}

func (c *fakeConn) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return c.rows, nil
}

func TestDriverRequiresCallbackConnection(t *testing.T) {
	d := New()
	if err := d.Exec(context.Background(), "SELECT 1", []any{}, nil); !errors.Is(err, ErrMissingConn) {
		t.Fatalf("Exec error = %v, want %v", err, ErrMissingConn)
	}
}

func TestDriverExecUsesContextConnection(t *testing.T) {
	d := New()
	conn := &fakeConn{tag: pgconn.NewCommandTag("INSERT 0 1")}
	ctx := WithConn(context.Background(), conn)

	var result entsql.Result
	if err := d.Exec(ctx, "INSERT INTO test DEFAULT VALUES", []any{}, &result); err != nil {
		t.Fatal(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatalf("RowsAffected = %d, want 1", affected)
	}
	if _, err := result.LastInsertId(); err == nil {
		t.Fatal("LastInsertId should be unsupported")
	}
}

func TestDriverQueryAdaptsColumnNames(t *testing.T) {
	d := New()
	rows := &fakeRows{fields: []pgconn.FieldDescription{{Name: "id"}, {Name: "name"}}}
	ctx := WithConn(context.Background(), &fakeConn{rows: rows})

	var result entsql.Rows
	if err := d.Query(ctx, "SELECT id, name FROM test", []any{}, &result); err != nil {
		t.Fatal(err)
	}
	columns, err := result.Columns()
	if err != nil {
		t.Fatal(err)
	}
	if len(columns) != 2 || columns[0] != "id" || columns[1] != "name" {
		t.Fatalf("Columns = %#v", columns)
	}
	if result.NextResultSet() {
		t.Fatal("pgx adapter must not report multiple result sets")
	}
	result.Close()
}

func TestDriverDisablesEntTransactions(t *testing.T) {
	if _, err := New().Tx(context.Background()); !errors.Is(err, ErrEntTransactionsDisabled) {
		t.Fatalf("Tx error = %v, want %v", err, ErrEntTransactionsDisabled)
	}
}

func TestDriverDialectAndClose(t *testing.T) {
	d := New()
	if got := d.Dialect(); got != dialect.Postgres {
		t.Fatalf("Dialect = %q, want %q", got, dialect.Postgres)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
}

type fakeRows struct {
	pgx.Rows
	fields []pgconn.FieldDescription
}

func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription { return r.fields }
func (r *fakeRows) Close()                                       {}
func (r *fakeRows) Err() error                                   { return nil }
func (r *fakeRows) Next() bool                                   { return false }
func (r *fakeRows) Scan(...any) error                            { return nil }
