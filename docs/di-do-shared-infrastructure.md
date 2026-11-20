# Shared infrastructure DI packages

The `samber/do/v2` migration starts by exposing infrastructure registrations as package-level `do.Package` `func(do.Injector)` values. Domain and application constructors continue to receive their dependencies explicitly; the injector stays at composition boundaries.

## Migrated in this slice

- `internal/lib/lf.Package` lazily provides the shared `*slog.Logger`.
- `internal/lib/kv.Package` lazily provides the `kv.Loader` port. It resolves the logger from the same injector and returns reference-configuration read errors to the caller instead of panicking.
- `internal/lib/db.Package` lazily provides PostgreSQL storage configuration and the `tx.Transactor`. Pool readiness is registered with `app.Lifecycle`, while pool closure is owned by a `samber/do` shutdown resource.
- `internal/lib/wp.Package` lazily provides both worker-pool implementations and registers shutdown resources. Application shutdown cancels the shared task context, stops accepting work, and waits for accepted tasks to finish, subject to the caller's shutdown context.
- `internal/lib/ws.Package` lazily provides the Echo server. Its lifecycle hook binds the TCP listener synchronously during startup, reports bind failures to `App.Start`, and gracefully shuts down the server.

A minimal composition is:

```go
injector := do.New(
    lf.Package,
    kv.Package,
)
loader, err := do.Invoke[kv.Loader](injector)
```

The configuration loader is lazy: the reference file is read when `kv.Loader` is first invoked. The required `reference.yaml` remains an eager startup requirement at the point of resolution; the optional `application.yaml` continues to fall back to defaults when absent.

## Intentionally not migrated in this slice

- The production entrypoint still uses the legacy Fx modules. The `samber/do` packages are available for the future composition-root switch; this slice does not change `cmd/engine/main.go` or remove Fx.

The legacy Fx modules remain available for the current production composition. This slice does not switch the production entrypoint or remove Fx.
