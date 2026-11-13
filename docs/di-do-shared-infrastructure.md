# Shared infrastructure DI packages

The `samber/do/v2` migration starts by exposing infrastructure registrations as package-level `do.Package` `func(do.Injector)` values. Domain and application constructors continue to receive their dependencies explicitly; the injector stays at composition boundaries.

## Migrated in this slice

- `internal/lib/lf.Package` lazily provides the shared `*slog.Logger`.
- `internal/lib/kv.Package` lazily provides the `kv.Loader` port. It resolves the logger from the same injector and returns reference-configuration read errors to the caller instead of panicking.

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

- `db` owns a PostgreSQL pool and currently couples pool ping/close to Fx hooks. Its `OnStop` callback also returns before `pool.Close` completes. It needs an explicit shutdown contract before registration moves.
- `ws` owns the Echo server and currently starts it asynchronously from an Fx hook without waiting for listener readiness. It needs explicit start/stop orchestration and startup error propagation.
- `wp` owns worker pools; their stop/wait semantics must be connected to the future App lifecycle rather than left to container garbage collection.

These providers remain on the existing Fx production composition until the lifecycle work establishes a single, explicit owner. This slice does not switch the production entrypoint or remove Fx.
