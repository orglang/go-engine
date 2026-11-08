// Package tx defines the public transaction and connection contracts.
package tx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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

// Transactor owns the lifetime of all resources passed to callbacks.
//
// ExplicitTx starts a PostgreSQL transaction and commits it when the callback
// succeeds; it rolls back when the callback fails. ImplicitTx acquires a
// connection without starting a transaction. Both callbacks receive the caller's
// context and a Conn that is valid only until the callback returns.
//
// Nested transactor calls are not supported. Conn deliberately exposes no
// transaction-management methods, so repository code cannot start transactions
// through the connection. Operations involving multiple drivers are not atomic
// across drivers.
type Transactor interface {
	ExplicitTx(context.Context, func(context.Context, Conn) error) error
	ImplicitTx(context.Context, func(context.Context, Conn) error) error
}
