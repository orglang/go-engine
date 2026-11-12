# Fx Composition Inventory

> Snapshot: `main` at `a6dd30971cb01996c5f7e36fbcde03c5ffdebf77` (2026-11-12).
> Scope: inventory only. No DI implementation or runtime behavior is changed by this document.

## 1. Repository and composition boundaries

The application source belongs to `orglang/go-engine`, not the `go-workspace` meta-repository. The workspace's `go.work` includes `./engine` and `./sdk`; each is a separate repository and Go module.

Current production entrypoint: `cmd/engine/main.go`.

It creates one Fx application directly in `main`:

- registers shared infrastructure modules;
- registers pool and proc adapter modules individually;
- registers `internal/app/web`;
- calls `.Run()`.

There is no explicit `App` type or separate `pool` / `proc` / `prog` composition function in the inspected production entrypoint. The `internal/prog` package has core and adapter code, but no Fx module and is not wired by `cmd/engine/main.go`.

## 2. Declared Fx modules

There are 21 files named `di_fx.go` in the engine tree.

| Area | Fx module | Declared in | Registered by production entrypoint |
|---|---|---|---|
| Shared infrastructure | `lib/db` | `internal/lib/db/di_fx.go` | Yes |
| Shared infrastructure | `lib/e2e` | `internal/lib/e2e/di_fx.go` | No — test-only composition |
| Shared infrastructure | `lib/kv` | `internal/lib/kv/di_fx.go` | Yes |
| Shared infrastructure | `lib/lf` | `internal/lib/lf/di_fx.go` | Yes |
| Shared infrastructure | `lib/wp` | `internal/lib/wp/di_fx.go` | Yes |
| Shared infrastructure | `lib/ws` | `internal/lib/ws/di_fx.go` | Yes |
| Pool | `pool/compconn` | `internal/pool/compconn/adapter/di_fx.go` | Yes |
| Pool | `pool/compexec` | `internal/pool/compexec/adapter/di_fx.go` | Yes |
| Pool | `pool/compvar` | `internal/pool/compvar/adapter/di_fx.go` | Yes |
| Pool | `pool/termdef` | `internal/pool/termdef/adapter/di_fx.go` | Yes |
| Pool | `pool/termexp` | `internal/pool/termexp/adapter/di_fx.go` | **No** |
| Pool | `pool/typedef` | `internal/pool/typedef/adapter/di_fx.go` | Yes |
| Pool | `pool/typeexp` | `internal/pool/typeexp/adapter/di_fx.go` | Yes |
| Proc | `proc/compconn` | `internal/proc/compconn/adapter/di_fx.go` | Yes |
| Proc | `proc/compexec` | `internal/proc/compexec/adapter/di_fx.go` | Yes |
| Proc | `proc/termdec` | `internal/proc/termdec/adapter/di_fx.go` | Yes |
| Proc | `proc/termdef` | `internal/proc/termdef/adapter/di_fx.go` | Yes |
| Proc | `proc/termexp` | `internal/proc/termexp/adapter/di_fx.go` | **No** |
| Proc | `proc/typedef` | `internal/proc/typedef/adapter/di_fx.go` | Yes |
| Proc | `proc/typeexp` | `internal/proc/typeexp/adapter/di_fx.go` | Yes |
| Inbound adapter | `app/web` | `internal/app/web/di_fx.go` | Yes |

Thus, 18 of the 20 production-oriented Fx modules are passed to the production root. The remaining declared module, `lib/e2e`, is deliberately used by a separate test composition (see section 5).

## 3. Seven production areas not currently wired by the production root

The following seven areas match the previously identified inventory gap. “Not wired” is a factual description of the inspected entrypoint, not a recommendation to register every item.

| Area | Existing code | Fx module | Initial observation / decision needed |
|---|---|---|---|
| `pool/termexp` | Core and adapter; adapter has `di_fx.go` | Declared, not passed to `fx.New` | Determine whether it is needed by a production use case or intentionally unused |
| `proc/termexp` | Core and adapter; adapter has `di_fx.go` | Declared, not passed to `fx.New` | Determine whether it is needed by a production use case or intentionally unused |
| `pool/compstep` | Core and adapter directories | No `di_fx.go` | Determine whether it needs providers or is a data/contract-only package |
| `proc/compstep` | Core and adapter directories | No `di_fx.go` | Determine whether it needs providers or is a data/contract-only package |
| `pool/connact` | Core and adapter directories | No `di_fx.go` | Determine whether it needs providers or is a data/contract-only package |
| `proc/connact` | Core and adapter directories | No `di_fx.go` | Determine whether it needs providers or is a data/contract-only package |
| `prog` | Core and adapter directories | No `di_fx.go` | Top-level program composition/API is not wired by the current production root |

For each area, the migration must record one explicit outcome:

1. wire it at the appropriate composition boundary, with a test proving the intended registration/resolution; or
2. document why it is intentionally not wired and what code path (if any) owns its use.

Do not create empty DI modules solely to make the inventory look complete. In particular, a package that only defines request/step data types may not need container registration.

## 4. Fx API usage found

| Fx API / feature | Observed usage |
|---|---|
| `fx.New` | Production entrypoint and E2E test setup |
| `fx.Module` | 21 declared `di_fx.go` modules |
| `fx.Provide` | Providers in the declared modules |
| `fx.Private` | Used to keep implementation-only providers private in several modules |
| `fx.Annotate` / `fx.As` | Used to expose concrete constructors as port/interface types |
| `fx.Invoke` | Used for HTTP route/controller/presenter and broker configuration |
| `fx.In` | Used by `pool/compexec` service parameter struct and E2E test suite |
| `fx.Out` | No usage found |
| `fx.Populate` | Used by E2E test setup to populate its suite |
| `fx.Lifecycle` / `fx.Hook` | Used by DB pool/transactor and Echo server constructors |
| `fx.From` | Only appears in commented-out code in `pool/compexec` wiring |
| `App.Run` / `App.Stop` | Production calls `Run`; E2E cleanup calls `Stop` |

### Lifecycle-sensitive constructors

- `internal/lib/db/core.go:newTransactor` receives `fx.Lifecycle` and passes it to `newPoolPgx`.
- `internal/lib/db/ds_pgx.go:newPoolPgx` registers `pool.Ping` as `OnStart` and a close callback as `OnStop`. The current callback starts `pool.Close()` in a goroutine and returns immediately; this needs an explicit shutdown-semantics review during the lifecycle migration.
- `internal/lib/ws/echo.go:newEchoServer` registers hooks that start Echo in a goroutine and call `e.Shutdown(ctx)` on stop. The start hook returns before the listener reports readiness; startup/error behavior needs to be preserved or deliberately improved.
- `cmd/engine/main.go` delegates process lifecycle to `fx.New(...).Run()`.

### Provider and wiring observations

- Most modules bind concrete adapters to core ports using `fx.Annotate(..., fx.As(...))`.
- `pool/compexec` has the broadest visible dependency set: its `serviceParams` uses `fx.In` and contains repositories, broker, transaction, and logger dependencies.
- `fx.Invoke` is not only dependency construction: it performs side-effectful route and broker configuration. These calls must be accounted for explicitly when replacing Fx.
- `pool/compexec` registers both a concrete `newPondBroker` and an interface binding for `compexec.Broker`; it also contains commented-out alternative wiring. Review for duplicate construction/ownership before migration.
- Some providers are constructed with arguments at registration time (for example, ADT DAO constructors receiving table metadata), rather than being passed as plain zero-argument provider functions.

## 5. Test composition and baseline

The E2E setup in `test/e2e/pool_test.go` builds a separate Fx application:

```go
fx.New(rc.Module, e2e.Module, fx.Populate(s))
```

It obtains six E2E APIs from `internal/lib/e2e`, opens a PostgreSQL connection directly for test assertions, and stops the Fx app from test cleanup. This is a separate composition root and must be included in the DI migration; it is not a missing production registration.

The engine Taskfile defines these checks:

| Task | Command | Notes |
|---|---|---|
| Unit | `go test -tags=unit ./internal/lib/... ./internal/adt/... ./internal/app/... ./cmd/... ./internal/pool/... ./internal/proc/... ./internal/prog/... ./internal/proto/...` | Main fast test baseline |
| Integration | `go test -run=TestIntegration ./internal/lib/... ./internal/adt/... ./internal/app/... ./cmd/...` | May require local infrastructure/configuration |
| E2E | `go test -failfast ./test/...` | Requires the test environment, including PostgreSQL |
| Sources check | `task sources:check` | Runs source generation, formatting diff, and golangci-lint |

Baseline evidence available from GitHub at inventory time:

- `main` HEAD is `a6dd30971cb01996c5f7e36fbcde03c5ffdebf77`.
- The latest `Sources Check` workflow run for that HEAD completed successfully on 2026-11-09.
- This inventory was performed through repository source inspection; the Go test commands above were **not executed locally as part of this change**. A fresh run should be recorded before changing composition.

## 6. Initial migration implications

1. Keep `lib/e2e` as a distinct test composition unless later evidence shows it should share production composition.
2. Treat the two declared but unregistered `termexp` modules separately from the five areas with no Fx module at all.
3. Move route/broker configuration currently hidden in `fx.Invoke` into explicit composition steps.
4. Replace constructor-level `fx.Lifecycle` dependencies with lifecycle ownership that can be tested independently of the DI container.
5. Inspect the concrete provider graph before deciding package boundaries; the current module-per-capability layout is not necessarily the desired `do.Package` layout.
6. Preserve the existing `tx.Conn` contract and transaction ownership model.

## 7. Scope limitations

This document does not claim that all possible dynamic or indirect Fx usage has been proven absent. The inventory covers the checked-in Go source and the current production/test composition entrypoints at the snapshot commit above. The SDK repository could not be inspected from the workspace's referenced repository name during this pass; no SDK source changes are part of this stage.
