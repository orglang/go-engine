# Pool domain DI composition

The `internal/pool` package is the composition boundary for pool-domain hexagons.
Its `Package` is intentionally migrated incrementally; it does not yet represent
the complete pool runtime.

## Migrated persistence slice

The first slice registers these repository ports through `samber/do` v2:

- `compconn.Repo` with its private SQL query builder.
- `compvar.Repo` with its private SQL query builder.
- `typeexp.Repo` with its private SQL query builder.

Each adapter package owns its concrete construction details and exposes only
the core repository and query-builder interfaces. Repositories continue to
receive `tx.Conn` from their callers; no provider acquires a connection or
changes transaction ownership.

## Not included yet

The remaining pool adapters (`compexec`, `termdef`, `typedef`) require additional
composition decisions: service dependencies, ADT repositories, broker
subscription, and HTTP route registration. They remain on the legacy Fx
composition until those concerns are moved into explicit, testable wiring.
`termexp` remains intentionally unwired because its repository has no production
consumer. Translation/model-only packages do not receive DI modules.

The production entrypoint remains on Fx during this incremental migration.
