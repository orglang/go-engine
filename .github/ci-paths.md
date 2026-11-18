# CI path classification

The change detector classifies pull-request paths on two independent axes:

- **Area:** `app` or `tool`.
- **Kind:** `deps`, `specs`, `sources`, `binaries`, or `distros`.

The resulting outputs are:

| Area | Dependency inputs | Specifications | Sources | Binaries | Distributions |
|---|---|---|---|---|---|
| Application | `app-deps` | `app-specs` | `app-sources` | `app-binaries` | `app-distros` |
| Tooling | `tool-deps` | `tool-specs` | `tool-sources` | `tool-binaries` | `tool-distros` |

## Implementation

The composite action at `.github/actions/ci-changes` uses `dorny/paths-filter@v4`, following the established pattern in `orglang/go-workspace`.

The filter definitions live in `.github/ci-paths.yaml`. Each filter includes the shared broad-invalidation paths so workflow/action definitions, Taskfiles, and mise configuration/lockfiles conservatively select all ten outputs.

The action only invokes the path filter for `pull_request` events. For every other event it reports `true` for all outputs, preserving full validation for push, merge queue, and manual runs. Pull request workflows using the action must grant `pull-requests: read`.

## Current path rules

- `app-deps`: Go module manifests and sums (`go.mod`, `go.sum`).
- `tool-deps`: TypeSpec npm manifests/lockfiles and mise toolchain manifests/lockfiles.
- `app-specs`: TypeSpec source files and generated OpenAPI.
- `tool-specs`: TypeSpec compiler configuration.
- `app-sources`: Go and SQL source files.
- `tool-sources`: repository scripts used by source/tooling workflows.
- `app-binaries`: application command entrypoint/build inputs under `cmd/engine/`.
- `app-distros`: Dockerfiles, Compose files, and Docker ignore configuration.
- `tool-distros`: distribution-specific mise manifests/lockfiles.
- `tool-binaries`: reserved for tool-specific binary build inputs; currently selected through broad tooling invalidation.

A path may select more than one group. For example, the application command entrypoint is both a Go source and an application binary input.

## Classification versus impact

These outputs describe which input groups changed. They are not a one-to-one list of jobs to run. Workflow jobs must apply a separate impact/fan-out mapping for upstream dependencies. For example, application dependencies affect binary and distribution checks even though the changed path is classified as `app-deps`; TypeSpec changes can affect generated sources and downstream artifacts.

If the impact of a path cannot be established safely, broaden invalidation. Do not add workflow-level path filters to required workflows.

## Delivery

This PR introduces and exercises the classification contract only. Existing CI jobs are intentionally not gated in this step.
