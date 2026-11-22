# Pool domain DI composition

The `internal/pool` package is the composition boundary for pool-domain hexagons.
Its `Package` is intentionally migrated incrementally; it does not yet represent
the complete pool runtime.

## Migrated persistence and service slice

The following repository ports are registered through `samber/do` v2:

- `compconn.Repo` backed by Ent and specialized pgx operations.
- `compvar.Repo` with its current persistence adapter.
- `compexec.Repo` backed by Ent and specialized pgx queries.
- `typeexp.Repo` backed by Ent.
- `typedef.Repo` backed by Ent.
- `termdef.Repo` backed by Ent.
- The shared `descsem.Repo` used by the definition and execution services.
- `compexec.API` with its process execution and ADT repository dependencies.
- The asynchronous `compexec.Broker` backed by the shared pond worker pool.

The process computation-execution repository has its own adapter package.
Each adapter package owns its concrete persistence details and exposes core
ports. Repositories continue to receive `tx.Conn` from their callers; no
provider acquires a connection or changes transaction ownership.

The broker and API form a runtime relationship, so subscription is not hidden
inside a lazy provider. The application composition root must call
`compexecadapter.SubscribeBroker(injector)` after registering `pool.Package`.

## HTTP boundary

The pool computation-execution adapter exposes
`compexecadapter.RegisterRoutes(echoServer, api, logger)` as an explicit
inbound composition function. The legacy Fx module invokes this function, so
the existing production routes remain active while the application entrypoint
is migrated in a later step. A focused adapter test verifies all three routes.

The production entrypoint remains on Fx during this incremental migration.

`termexp` remains intentionally unwired because its repository has no
production consumer. Translation/model-only packages do not receive DI modules.
