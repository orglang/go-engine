# Pool domain DI composition

The `internal/pool` package is the composition boundary for pool-domain hexagons.
Its `Package` is intentionally migrated incrementally; it does not yet represent
the complete pool runtime.

## Migrated persistence and service slice

The following repository ports are registered through `samber/do` v2:

- `compconn.Repo` with its private SQL query builder.
- `compvar.Repo` with its private SQL query builder.
- `compexec.Repo` with its private SQL query builder.
- `typeexp.Repo` with its private SQL query builder.
- `typedef.Repo` with its private SQL query builder.
- `termdef.Repo` with its private SQL query builder.
- The shared `descsem.Repo` used by the definition and execution services.
- `compexec.API` with its process execution and ADT repository dependencies.
- The asynchronous `compexec.Broker` backed by the shared pond worker pool.

The process computation-execution repository has its own adapter package and
query builder. Each adapter package owns its concrete construction details and
exposes core ports. Repositories continue to receive `tx.Conn` from their
callers; no provider acquires a connection or changes transaction ownership.

The broker and API form a runtime relationship, so subscription is not hidden
inside a lazy provider. The application composition root must call
`compexecadapter.SubscribeBroker(injector)` after registering `pool.Package`.

## Not included yet

HTTP controller construction and route registration remain on the legacy Fx
composition and will move at the inbound/composition boundary. The production
entrypoint remains on Fx during this incremental migration.

`termexp` remains intentionally unwired because its repository has no
production consumer. Translation/model-only packages do not receive DI modules.
