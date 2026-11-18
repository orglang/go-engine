# Change detection

The change detector classifies pull-request paths by the **earliest pipeline stage that consumes them**. Classification describes what changed; downstream impact determines which later stages must also run.

## Input groups

There are four public stage groups. Within each group, comments distinguish application inputs, dependencies, and tooling where those categories exist.

| Stage | Application inputs | Dependencies | Tooling |
|---|---|---|---|
| specs | TypeSpec definitions and generated OpenAPI contract | TypeSpec npm manifests | TypeSpec toolchain |
| sources | Go source code | — | Sources toolchain |
| binaries | Database bootstrap and migrations | Go module manifests (go.mod, go.sum) | — |
| distros | Dockerfiles and Compose definitions | — | Distribution-specific toolchain, when defined here |

Shared orchestration/configuration is an internal catch-all filter, not a public stage group:

- GitHub Actions workflows/actions and the change-group definition.
- Taskfiles.
- Root mise configuration and lockfile, plus .miserc.toml.

A match in the `any` filter conservatively marks all four public stage outputs as changed. Stage-specific mise toolchains are classified with the earliest stage that consumes them: mise.specs-dev.* under specs, and mise.sources-dev.* under sources.

## Impact mapping

A changed stage input invalidates that stage and every downstream stage:

- specs invalidates specs, sources, binaries, and distros.
- sources invalidates sources, binaries, and distros.
- binaries invalidates binaries and distros.
- distros invalidates only distros.

The `any` filter invalidates all four stages. When a narrower impact cannot be proven safe, prefer broader invalidation.

Do not use workflow-level path filters on required workflows: skipped workflows can leave required checks pending. Workflow jobs should use this contract to decide whether to run while allowing the required workflows themselves to start.

## Dependency consistency checks

Dependency-lock consistency is checked by the stage that owns each dependency set:

- specs:check refreshes the TypeSpec toolchain lock and npm package lock, then verifies neither changed.
- sources:check refreshes the Sources toolchain lock, then verifies it and generated Go files are unchanged.
- binaries:check runs GOWORK=off go mod tidy and mise lock, then verifies go.mod, go.sum, and mise.lock are unchanged.

The Generation maintenance workflow at `.github/workflows/generation.yaml` refreshes all three sets explicitly by running `task specs:lock`, `task sources:lock`, and `task binaries:lock` before generation. There is no aggregate `deps:lock` task; keeping the stage ownership explicit avoids a second, hidden orchestration entry point.

Unit tests run in Binaries Build alongside integration tests and binary compilation. This ensures Go module changes are validated without requiring Sources Verify to run.

## Action

The composite action at .github/actions/change-detection uses dorny/paths-filter@v4 and reads .github/change-groups.yaml.

It exposes four stage outputs: specs, sources, binaries, and distros. It runs path filtering only for pull_request; other events report true for every output so merge-queue and manual runs retain full validation. Pull-request workflows using the action must grant pull-requests: read.

Specs Check uses the `specs` output to skip its expensive job when a pull request does not change specs inputs. Sources Verify uses `specs` and `sources`: a specs change also invalidates source verification under the downstream impact mapping. Binaries Build uses `specs`, `sources`, and `binaries`; Distro Build uses all four outputs because every earlier stage can invalidate distribution artifacts. Each job also skips draft pull requests. Merge-group and workflow_dispatch runs continue to execute all four checks. Push runs continue to execute the workflows configured for push events.
