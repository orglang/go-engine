# Pool domain DI composition

The `internal/pool` package is the composition boundary for pool-domain hexagons.
Its `Package` is intentionally migrated incrementally; it does not yet represent
the complete pool runtime.

## Migrated persistence and service slice

The following repository ports are registered through `samber/do` v2:

- `execconn.Repo` backed by Ent and specialized pgx operations.
- `execvar.Repo` with its current persistence adapter.
- `execeval.Repo` backed by Ent and specialized pgx queries.
- `typeexp.Repo` backed by Ent.
- `typedef.Repo` backed by Ent.
- `termdef.Repo` backed by Ent.
- Description-qualified names are stored by the existing `typedef.Repo` and `termdef.Repo` adapters.
- `execeval.API` with its process execution and ADT repository dependencies.
- The asynchronous `execeval.Broker` backed by the shared pond worker pool.

The process computation-execution repository has its own adapter package.
Each adapter package owns its concrete persistence details and exposes core
ports. Repositories continue to receive `tx.Conn` from their callers; no
provider acquires a connection or changes transaction ownership.

The broker and API form a runtime relationship, so subscription is not hidden
inside a lazy provider. The application composition root must call
`execevaladapter.SubscribeBroker(injector)` after registering `pool.Package`.

## HTTP boundary

The pool computation-execution adapter exposes
`execevaladapter.RegisterRoutes(echoServer, api, logger)` as an explicit
inbound composition function. The legacy Fx module invokes this function, so
the existing production routes remain active while the application entrypoint
is migrated in a later step. A focused adapter test verifies all three routes.

The production entrypoint remains on Fx during this incremental migration.

Runtime expression DTOs and their bidirectional mapping live in `execeval/adapter`;
the former `termexp` application package has been removed. Translation/model-only
packages do not receive DI modules.
