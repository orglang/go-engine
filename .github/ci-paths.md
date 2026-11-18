# CI path classification

The change detector classifies pull-request paths on two independent axes:

- **Area:** `app` or `tool`.
- **Kind:** `deps`, `specs`, `sources`, `binaries`, or `distros`.

The resulting outputs are:

| Area | Dependency inputs | Specifications | Sources | Binaries | Distributions |
|---|---|---|---|---|---|
| Application | `app-deps` | `app-specs` | `app-sources` | `app-binaries` | `app-distros` |
| Tooling | `tool-deps` | `tool-specs` | `tool-sources` | `tool-binaries` | `tool-distros` |

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

Changes to workflow/action definitions, Taskfiles, and mise configuration/lockfiles conservatively select all ten outputs. These files can change how multiple checks are built or executed, and the initial detector deliberately prefers false positives over missed checks.

A path may select more than one group. For example, the application command entrypoint is both a Go source and an application binary input.

## Classification versus impact

These outputs describe which input groups changed. They are not a one-to-one list of jobs to run. Workflow jobs must apply a separate impact/fan-out mapping for upstream dependencies. For example, application dependencies affect binary and distribution checks even though the changed path is classified as `app-deps`; TypeSpec changes can affect generated sources and downstream artifacts.

If the impact of a path cannot be established safely, broaden invalidation. Do not add workflow-level path filters to required workflows.

## Using the action

The composite action at `.github/actions/ci-changes` expects a checkout with sufficient history for the three-dot diff. For pull requests, pass the base and head SHAs from the event and use `fetch-depth: 0`. For every non-`pull_request` event it returns `true` for all outputs, preserving merge-queue, push, and manual-run behavior.

The first implementation introduces and tests the contract only. Existing workflows are intentionally not gated in this step.
