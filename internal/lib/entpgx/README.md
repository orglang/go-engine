# Context-aware Ent driver for pgx

This package is an infrastructure spike for using one application-scoped
`ent.Client` with the pgx resource supplied to each `tx.Transactor` callback.

## Usage

Construct one Ent client in dependency injection:

```go
client := ent.NewClient(ent.Driver(entpgx.New()))
```

Inside a repository adapter, bind the callback resource to the context used for
all Ent operations:

```go
entCtx := entpgx.WithConn(ctx, conn.Pgx)
err := client.User.Create().SetName(name).Exec(entCtx)
```

The driver is stateless. It never stores a connection or transaction, and
routes each operation through the resource attached to that operation's
context. The context must not escape the transactor callback.

Native pgx operations continue to use `conn.Pgx`, so Ent and pgx statements
inside one callback use the same underlying `pgx.Tx` (for `ExplicitTx`) or
acquired `*pgxpool.Conn` (for `ImplicitTx`).

## Deliberate restrictions

- Ent-managed transactions (`client.Tx`) are disabled. The application's
  `tx.Transactor` exclusively owns begin/commit/rollback.
- `Close` is a no-op. Pool lifecycle belongs to the existing DB lifecycle.
- Calls without a context-bound pgx resource fail with `ErrMissingConn`.
- PostgreSQL `LastInsertId` is unsupported.
- `ColumnTypes` returns an empty slice because Ent's generated entity
  scanners need column names and scan destinations, while pgx exposes OIDs and
  `database/sql` has no public constructor for `sql.ColumnType`.
- pgx does not expose multiple result sets; `NextResultSet` always returns
  false.

## Required verification before adoption

This is not yet a production-approved adapter. Run unit tests and PostgreSQL
integration tests against the repository's pinned Go toolchain and Ent version.
In particular verify:

1. Ent CRUD and `INSERT ... RETURNING`.
2. Nullable fields, `jsonb`, and `ltree` scan/encode behavior.
3. Mixed Ent + pgx writes commit atomically.
4. Mixed Ent + pgx writes roll back atomically on callback error.
5. Two concurrent callbacks never cross-route their pgx resources.
6. Missing context resource and Ent-managed transaction attempts fail clearly.

Do not wire this driver into production adapters until these checks pass.
