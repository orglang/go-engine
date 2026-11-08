// Package tx defines the public transaction and connection contracts.
package tx

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrMissingPgx is returned when a connection does not contain a PostgreSQL resource.
var ErrMissingPgx = errors.New("connection does not contain a pgx resource")

// Pgx is the minimal set of PostgreSQL operations currently used by adapters.
// pgx.Tx and *pgxpool.Conn satisfy this interface.
type Pgx interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

// Conn contains driver resources for one callback invocation.
// It must not be retained or used after the callback returns.
type Conn struct {
	Pgx Pgx
}

// RequirePgx returns the PostgreSQL resource or a regular error when it is absent.
func (c Conn) RequirePgx() (Pgx, error) {
	if c.Pgx == nil {
		return nil, ErrMissingPgx
	}
	return c.Pgx, nil
}

// UoW is the legacy callback payload. Use Conn and the callback context in new code.
//
// Deprecated: use Conn with Transactor.ExplicitTx or Transactor.ImplicitTx.
type UoW struct {
	Ctx context.Context
	Pgx Pgx
}

// Transactor owns the lifetime of all resources passed to callbacks.
//
// ExplicitTx starts a PostgreSQL transaction and commits it when the callback
// succeeds; it rolls back when the callback fails. ImplicitTx acquires a
// connection without starting a transaction. Both callbacks receive a derived
// context and a Conn that is valid only until the callback returns.
//
// Nested transactor calls are not supported. Conn deliberately exposes no
// transaction-management methods, so repository code cannot start transactions
// through the connection. Operations involving multiple drivers are not atomic
// across drivers.
//
// The Old methods preserve the previous callback contract during migration.
type Transactor interface {
	ExplicitTx(context.Context, func(context.Context, Conn) error) error
	ImplicitTx(context.Context, func(context.Context, Conn) error) error

	// Deprecated: use ExplicitTx.
	ExplicitTxOld(context.Context, func(UoW) error) error
	// Deprecated: use ImplicitTx.
	ImplicitTxOld(context.Context, func(UoW) error) error
}
