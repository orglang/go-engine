# Change detection

The change detector classifies pull-request paths by the **earliest pipeline stage that consumes them**. Classification describes what changed; downstream impact determines which later stages must also run.

## Input groups

There are four public stage groups. Within each group, comments distinguish application inputs, dependencies, and tooling where those categories exist.

| Stage | Application inputs | Dependencies | Tooling |
|---|---|---|---|
| specs | TypeSpec definitions and generated OpenAPI contract | TypeSpec npm manifests | TypeSpec toolchain |
| sources | Go source code | — | Sources toolchain and lint configuration |
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

Unit tests run in Binaries Gate alongside integration tests and binary compilation. This ensures Go module changes are validated without requiring Sources Gate to run.

## Action

The composite action at .github/actions/change-detection uses dorny/paths-filter@v4 and reads .github/change-groups.yaml.

It exposes four stage outputs: specs, sources, binaries, and distros. It runs path filtering only for pull_request; other events report true for every output so merge-queue and manual runs retain full validation. Pull-request workflows using the action must grant pull-requests: read.

Specs Gate uses the `specs` output to skip its expensive job when a pull request does not change specs inputs. Sources Gate uses `specs` and `sources`: a specs change also invalidates source verification under the downstream impact mapping. Binaries Gate uses `specs`, `sources`, and `binaries`; Distros Gate uses all four outputs because every earlier stage can invalidate distribution artifacts. The workflow display names are `Specs`, `Sources`, `Binaries`, and `Distros`; their primary job display names are `Specs Gate`, `Sources Gate`, `Binaries Gate`, and `Distros Gate`, with matching job IDs `specs-gate`, `sources-gate`, `binaries-gate`, and `distros-gate`. Each job also skips draft pull requests.

All four workflows run for merge-group and workflow_dispatch events. For push events, the detector selects every stage; only workflows configured with a push trigger start. Currently Specs Check and Sources Verify have push triggers, while Binaries Build and Distro Build do not. This distinction is intentional: change detection does not add or broaden workflow triggers.

## Task-level gates and artifact cache

CI keeps the path-based stage gates above. After a stage job starts, its mise gate selects the work inside that stage with `mise run --affected`:

- `specs:gate` selects specs checks for specs inputs.
- `sources:gate` selects source checks for specs or sources inputs.
- `binaries:gate` selects binary checks for specs, sources, or binaries inputs.
- `distros:gate` selects distribution checks for any engine stage input. It also performs the engine source preparation and workspace distribution checks, so those steps are skipped together when the gate is unaffected.

The `sources` declared by each gate mirror the downstream impact mapping above. Keep the two layers aligned: #87 decides whether a stage job starts; mise decides whether that stage's work is affected. Do not replace the existing job-level checks with mise gates.

Within `mise.toml`, the gates use named `task_config.input_groups` to compose cumulative stage inputs without repeating path lists. Keep the gates independent: each stage has its own GitHub Actions workflow, and a `depends` chain would execute upstream gates again inside downstream workflows. Use task dependencies for actual task prerequisites, not to encode CI stage ordering.

For pull requests, the gated jobs fetch full Git history and compare `HEAD` with `origin/$GITHUB_BASE_REF`. Merge-group, push, and manual runs execute the stage task unconditionally so they retain full validation.

The mise artifact cache is enabled only for `sources:lint` and `tests:unit`. Both declare `outputs = []`: they do not produce workspace files that later steps need. Do not apply this setting to generation or build tasks unless every required output is declared and can be restored correctly.

CI persists only `$MISE_CACHE_DIR/task-artifacts/v2` with GitHub Actions Cache. Each run writes to a unique key and restores from the OS-specific prefix; mise's own task keys decide whether an entry is reusable. GitHub's cache branch scoping keeps pull-request writes separate from the base branch cache.

## Experimental Pipe workflows

The `Specs Pipe`, `Sources Pipe`, `Binaries Pipe`, and `Distros Pipe` workflows are an opt-in comparison of cumulative stage execution. They use the same PR path-detection contract and stage admission conditions as the regular workflows, but run a separate mise task chain:

- `specs:pipe` runs Specs.
- `sources:pipe` depends on `specs:pipe` and then runs Sources.
- `binaries:pipe` depends on `sources:pipe` and then runs Binaries.
- `distros:pipe` depends on `binaries:pipe` and then runs Distros.

The regular `*:gate` tasks and their workflows remain independent and unchanged. Pipe tasks are selected with `mise run --affected`; once a Pipe task is selected, mise expands its dependencies, so upstream stages run as prerequisites even when their own inputs did not change. This deliberately measures the cost and behavior of cumulative execution against the independent-stage baseline.

Pipe workflows are not required checks. Their artifact cache uses a separate GitHub Actions cache namespace from the regular workflows to avoid the two variants warming each other's cache. Compare runs for the same PR revisions and inspect task timings/cache stats before drawing conclusions.

## Maintaining the mapping

When adding or moving an input path:

1. Add it to the earliest stage in `.github/change-groups.yaml` that consumes it. Keep dependency manifests with the stage that owns those dependencies; do not classify them as source code merely because they support a build.
2. If an input is shared orchestration or its impact cannot be bounded safely, add it to `any` rather than guessing a narrower stage.
3. Check the impact mapping above. A stage output must be included in that stage's job and every downstream job it invalidates.
4. If adding a workflow or a new gated job, use the shared composite action, expose only the outputs the job needs, grant `pull-requests: read` for PR path filtering, and preserve the required job's existing name.
5. Keep event triggers explicit. Non-PR events are deliberately unfiltered by the action; a workflow runs for an event only if that workflow declares the corresponding trigger.

Before merging mapping or gating changes, review representative isolated changes for each affected group, combinations of upstream and downstream groups, and shared/catch-all changes. Also verify that merge_group and workflow_dispatch select every stage, and that push-triggered workflows select every stage. Do not add workflow-level `paths` filters to required workflows.
