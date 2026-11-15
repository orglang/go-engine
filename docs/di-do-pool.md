# Pool domain DI composition

The `internal/pool` package is the composition boundary for pool-domain hexagons.
Its `Package` is intentionally migrated incrementally; it does not yet represent
the complete pool runtime.

## Migrated persistence slice

The following repository ports are registered through `samber/do` v2:

- `compconn.Repo` with its private SQL query builder.
- `compvar.Repo` with its private SQL query builder.
- `compexec.Repo` with its private SQL query builder.
- `typeexp.Repo` with its private SQL query builder.
- `typedef.Repo` with its private SQL query builder.
- `termdef.Repo` with its private SQL query builder.
- The shared `descsem.Repo` used by the definition services.

Each adapter package owns its concrete construction details and exposes only
the core repository and query-builder interfaces. Repositories continue to
receive `tx.Conn` from their callers; no provider acquires a connection or
changes transaction ownership.

## Not included yet

The computation-execution service and broker are not part of this slice.
Its service graph still needs the process execution repository and ADT
repositories, as well as an explicit broker subscription/lifecycle design.
HTTP controller construction and route registration also remain on the legacy
Fx composition. These concerns will be migrated as explicit composition
steps rather than hidden in provider side effects.

`termexp` remains intentionally unwired because its repository has no
production consumer. Translation/model-only packages do not receive DI modules.
The production entrypoint remains on Fx during this incremental migration.
