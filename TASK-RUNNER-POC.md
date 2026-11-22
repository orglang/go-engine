# Task runner migration PoC (issue #97)

This branch adds an isolated native-mise implementation for comparison. It does not replace or edit the existing Taskfile-backed tasks.

## Scope

- In scope: `go-engine` workflows and invoking them from the `go-workspace` root.
- Out of scope: `sdk`.
- Baseline: the existing `taskfile.yaml` and the current `mise.toml`.
- Experiment: `mise.poc.toml`, loaded as an additive mise environment named `poc`.

## Run from the engine repository

Run the existing Task implementation:

```sh
mise exec -- task specs:gate
mise exec -- task sources:gate
mise exec -- task binaries:build
mise exec -- go test -tags=unit ./internal/... ./cmd/...
mise exec -- go test -run=TestIntegration ./internal/... ./cmd/... ./test/it/...
mise exec -- task tests:e2e
```

Run the native-mise PoC:

```sh
mise -E poc run poc:specs:gate
mise -E poc run poc:sources:gate
mise -E poc run poc:binaries:build
mise -E poc run poc:binaries:unit-test
mise -E poc run poc:binaries:integration-test
mise -E poc run poc:tests:e2e
```

## Run engine tasks from the workspace root

From the `go-workspace` repository root (assuming the nested checkout is `engine/`):

```sh
mise -C engine -E poc run poc:binaries:build
mise -C engine -E poc run poc:binaries:unit-test
```

This deliberately exercises the same nested-repository boundary used by the workspace. No workspace config or `sdk` files are changed by this PoC.

## What is covered

| Workflow | Native PoC task | Notes |
|---|---|---|
| TypeSpec dependencies / generation / lint / analyze | `poc:specs:gate` | Explicit sequential composition; same npm and mise lock checks as Taskfile |
| Go source generation | `poc:sources:generate` | Same generator commands and ordering |
| Go source lint / analyze | `poc:sources:lint`, `poc:sources:analyze` | Same golangci-lint and go vet arguments |
| Sources gate | `poc:sources:gate` | Reproduces lock, generation, generated-file cleanliness, formatting, lint and analysis checks |
| Build / unit / integration / E2E | `poc:binaries:build`, `poc:binaries:unit-test`, `poc:binaries:integration-test`, `poc:tests:e2e` | Direct Go commands; avoids Task wrappers |

The PoC intentionally does not port the lock-refresh/check workflow (`binaries:lock` / `binaries:check`) or process/container convenience tasks yet. Those are useful follow-up cases if the first comparison indicates that native mise is viable.

## Evaluation checklist

- [ ] Both implementations run successfully on the same clean checkout.
- [ ] Generated files and lockfiles are unchanged after each gate.
- [ ] Failure exit codes and useful error output are equivalent.
- [ ] Tool versions and PATH are identical for each command.
- [ ] Workspace-root invocation resolves the engine config and working directory correctly.
- [ ] Compare task graph readability, parallelism, startup/runtime, and maintenance cost.
- [ ] Record any behavior that depends on Task-specific features (includes, aliases, vars, dirs, deps).

Do not remove Task or switch CI in this PoC. The outcome should be a decision backed by the same local runs, not a premature migration.
