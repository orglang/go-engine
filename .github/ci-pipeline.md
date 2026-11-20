# CI pipeline execution

CI workflows do not filter jobs by changed paths. When a workflow is triggered, its gate runs regardless of which files changed. This keeps required checks predictable and avoids a mismatch between changed-file detection and the dependency graph in mise.

## Engine stage pipeline

The engine mise tasks form a cumulative chain:

`specs:gate` → `sources:gate` → `binaries:gate`

Each downstream gate depends on the preceding gate rather than duplicating its inputs. The Binaries workflow invokes `mise run binaries:gate`, so Go source changes and dependency changes both execute the full prerequisite chain and the binary checks.

The Sources workflow invokes `mise run sources:gate`. The standalone Specs and Sources verification workflows also run their respective checks when triggered.

## Workspace distribution pipeline

The workspace repository owns `distros:gate`, because distribution builds and stack checks are workspace tasks. It depends on the engine's `//engine:binaries:gate`, then runs workspace `distros` and `stack commons` tasks in sequence.

## Workflow triggers

Each workflow declares its own triggers explicitly, including pull requests, merge groups, pushes, schedules, or manual dispatch where applicable. Pull-request jobs may still distinguish draft from ready-for-review PRs; that controls lifecycle, not changed-path filtering.

Do not add workflow-level `paths` filters to required workflows: skipped workflows can leave required checks pending. Keep stage ordering and dependency relationships in mise task definitions rather than reproducing them in GitHub Actions change-detection logic.
