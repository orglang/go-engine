# Change detection

The change detector classifies pull-request paths by the **earliest pipeline stage that consumes them**. Classification describes what changed; downstream impact determines which later stages must also run.

## Input groups

There are four public stage groups. Within each group, comments distinguish application inputs, dependencies, and tooling where those categories exist.

| Stage | Application inputs | Dependencies | Tooling |
|---|---|---|---|
| specs | TypeSpec definitions and generated OpenAPI contract | TypeSpec npm manifests | TypeSpec toolchain |
| sources | Go source code | — | Source generators, linters, and their toolchain |
| binaries | Database bootstrap and migrations | Go module manifests (go.mod, go.sum) | — |
| distros | Dockerfiles and Compose definitions | — | Distribution-specific toolchain, when defined here |

Shared orchestration/configuration is an internal filter, not a public stage group:

- GitHub Actions workflows/actions and the change-group definition.
- Taskfiles.
- Root mise configuration and lockfile, plus .miserc.toml.

A shared-configuration change conservatively marks all four public stage outputs as changed. Stage-specific mise toolchains are classified with the earliest stage that consumes them: mise.specs-dev.* under specs, and mise.sources-dev.* under sources.

## Impact mapping

A changed stage input invalidates that stage and every downstream stage:

- specs invalidates specs, sources, binaries, and distros.
- sources invalidates sources, binaries, and distros.
- binaries invalidates binaries and distros.
- distros invalidates only distros.

Shared configuration invalidates all four stages. When a narrower impact cannot be proven safe, prefer broader invalidation.

Do not use workflow-level path filters on required workflows: skipped workflows can leave required checks pending. Workflow jobs should use this contract to decide whether to run while allowing the required workflows themselves to start.

## Dependency consistency checks

Dependency-lock consistency is checked by the stage that owns each dependency set:

- specs:check refreshes the TypeSpec toolchain lock and npm package lock, then verifies neither changed.
- sources:check refreshes the root mise lock and Sources toolchain lock, then verifies neither changed.
- binaries:check runs GOWORK=off go mod tidy and verifies go.mod and go.sum are unchanged.

Unit tests run in Binaries Build alongside integration tests and binary compilation. This ensures Go module changes are validated without requiring Sources Verify to run.

The developer-facing deps:lock task remains available to update all dependency manifests and lockfiles together; the aggregate CI deps:check task is removed.

## Action

The composite action at .github/actions/change-detection uses dorny/paths-filter@v4 and reads .github/change-groups.yaml.

It exposes four stage outputs: specs, sources, binaries, and distros. It runs path filtering only for pull_request; other events report true for every output so push, merge-queue, and manual runs retain full validation. Pull-request workflows using the action must grant pull-requests: read.

This change introduces and exercises the classification contract only. Existing CI jobs are not gated yet.
