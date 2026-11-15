// Package entpgx provides a context-aware Ent driver backed by the pgx
// connection supplied by the application's transaction callback.
package entpgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"orglang/go-engine/internal/lib/tx"
)

var (
	// ErrMissingConn indicates that an Ent operation was invoked without a
	// connection supplied by the transaction callback.
	ErrMissingConn = errors.New("entpgx: pgx connection missing from context")

	// ErrEntTransactionsDisabled prevents Ent from creating a transaction
	// outside the application's Transactor lifecycle.
	ErrEntTransactionsDisabled = errors.New("entpgx: Ent-managed transactions are disabled; use tx.Transactor")
)

type connKey struct{}

// WithConn binds the callback-scoped pgx resource to ctx for Ent operations.
// The returned context must not outlive the Transactor callback.
func WithConn(ctx context.Context, conn tx.Pgx) context.Context {
	return context.WithValue(ctx, connKey{}, conn)
}

func connFromContext(ctx context.Context) (tx.Pgx, error) {
	conn, ok := ctx.Value(connKey{}).(tx.Pgx)
	if !ok || conn == nil {
		return nil, ErrMissingConn
	}
	return conn, nil
}

// Driver is stateless and safe for concurrent use. It routes each operation
// through the pgx resource attached to that operation's context.
type Driver struct{}

var _ dialect.Driver = (*Driver)(nil)

// New returns a stateless driver suitable for a single application-wide
// ent.Client. It does not own a pool, connection, or transaction.
func New() *Driver { return &Driver{} }

func (d *Driver) Dialect() string { return dialect.Postgres }

// Close is a no-op because connection and pool lifecycle belong to Transactor.
func (d *Driver) Close() error { return nil }

// Tx is deliberately disabled. The application Transactor is the only owner
// of transaction lifecycle; client.Tx would otherwise suggest nested ownership.
func (d *Driver) Tx(context.Context) (dialect.Tx, error) {
	return nil, ErrEntTransactionsDisabled
}

func (d *Driver) Exec(ctx context.Context, query string, args, result any) error {
	conn, err := connFromContext(ctx)
	if err != nil {
		return err
	}
	argv, ok := args.([]any)
	if !ok {
		return fmt.Errorf("entpgx: invalid args type %T, want []any", args)
	}

	switch result := result.(type) {
	case nil:
		_, err = conn.Exec(ctx, query, argv...)
		return err
	case *entsql.Result:
		tag, execErr := conn.Exec(ctx, query, argv...)
		if execErr != nil {
			return execErr
		}
		*result = commandTagResult{rowsAffected: tag.RowsAffected()}
		return nil
	default:
		return fmt.Errorf("entpgx: invalid Exec result type %T", result)
	}
}

func (d *Driver) Query(ctx context.Context, query string, args, result any) error {
	conn, err := connFromContext(ctx)
	if err != nil {
		return err
	}
	argv, ok := args.([]any)
	if !ok {
		return fmt.Errorf("entpgx: invalid args type %T, want []any", args)
	}
	rowsResult, ok := result.(*entsql.Rows)
	if !ok {
		return fmt.Errorf("entpgx: invalid Query result type %T, want *dialect/sql.Rows", result)
	}

	rows, err := conn.Query(ctx, query, argv...)
	if err != nil {
		return err
	}
	*rowsResult = entsql.Rows{ColumnScanner: &rowScanner{Rows: rows}}
	return nil
}

type commandTagResult struct {
	rowsAffected int64
}

func (r commandTagResult) LastInsertId() (int64, error) {
	return 0, errors.New("entpgx: LastInsertId is not supported by PostgreSQL")
}

func (r commandTagResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

var _ sql.Result = commandTagResult{}

// rowScanner adapts pgx.Rows to the subset of database/sql.Rows exposed by Ent.
type rowScanner struct {
	pgx.Rows
}

var _ entsql.ColumnScanner = (*rowScanner)(nil)

func (r *rowScanner) Columns() ([]string, error) {
	fields := r.FieldDescriptions()
	columns := make([]string, len(fields))
	for i := range fields {
		columns[i] = fields[i].Name
	}
	return columns, nil
}

// ColumnTypes metadata is not required by Ent's generated entity scanners.
// pgx exposes PostgreSQL OIDs rather than database/sql.ColumnType values, and
// database/sql does not provide a public constructor for ColumnType.
func (r *rowScanner) ColumnTypes() ([]*sql.ColumnType, error) {
	return []*sql.ColumnType{}, nil
}

// pgx.Rows does not expose multiple result sets.
func (r *rowScanner) NextResultSet() bool { return false }

func (r *rowScanner) Close() error {
	r.Rows.Close()
	return nil
}
