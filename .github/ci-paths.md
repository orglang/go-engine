# Change detection

The change detector classifies pull-request paths into independent input groups. Classification describes **what changed**; it does not decide by itself which CI jobs must run.

## Input groups

The application/tooling matrix has ten groups:

| Area | Dependencies | Specifications | Sources | Binaries | Distributions |
|---|---|---|---|---|---|
| Application | `app-deps` | `app-specs` | `app-sources` | `app-binaries` | `app-distros` |
| Tooling | `tool-deps` | `tool-specs` | `tool-sources` | `tool-binaries` | `tool-distros` |

Shared repository configuration is classified separately:

| Group | Inputs |
|---|---|
| `ci-config` | GitHub Actions workflows/actions and this change-group configuration |
| `task-config` | Taskfiles |
| `toolchain-config` | mise manifests, lockfiles, and mise configuration |

This keeps application dependencies (`go.mod`, `go.sum`) separate from tool dependencies (currently the TypeSpec npm manifest and lockfile). mise files are not package dependencies: they describe the toolchain. Taskfiles describe task orchestration. Neither is duplicated across all ten input groups.

A path can belong to more than one input group when it genuinely serves multiple roles (for example, a Go file under `cmd/engine/` is both application source and an application binary input).

## Impact mapping

The groups above are change classification, not a one-to-one job list. Workflows must apply a separate impact mapping:

- `app-deps` invalidates application binaries and downstream distributions.
- `tool-deps` invalidates checks that install or execute the affected tooling, plus their downstream outputs.
- `app-specs` invalidates specification checks and downstream generated sources, binaries, and distributions.
- `tool-specs` invalidates specification-generation/checking and downstream outputs.
- `app-sources` invalidates source checks, tests, binaries, and distributions as applicable.
- `tool-sources` invalidates checks that execute those scripts and their downstream outputs.
- `app-binaries` and `app-distros` invalidate their corresponding build/package checks.
- `tool-binaries` and `tool-distros` are reserved for future tool artifact pipelines.
- `ci-config` invalidates all ten checks because workflow/action behavior can change.
- `task-config` and `toolchain-config` currently invalidate all ten checks conservatively. This is an explicit fallback until task-to-workflow and toolchain-to-workflow dependencies are documented; it is not a claim that these files are inputs to every group.

When a narrower impact cannot be proven safe, prefer broader invalidation. Do not use workflow-level path filters on required workflows, because skipped workflows can leave required checks pending.

## Action

The composite action at `.github/actions/change-detection` uses `dorny/paths-filter@v4` and reads `.github/change-groups.yaml`.

It exposes the ten application/tooling outputs and the three shared-configuration outputs above. It runs path filtering only for `pull_request`; other events report `true` for every output so push, merge-queue, and manual runs retain full validation. Pull-request workflows using the action must grant `pull-requests: read`.

This change introduces and exercises the classification contract only. Existing CI jobs are not gated yet.
