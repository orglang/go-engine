# Go client generator comparison

Tracks PR 3 from [orglang/go-workspace#114](https://github.com/orglang/go-workspace/issues/114).

## Status

**Draft / evaluation in progress.** ogen is the leading candidate, but this document does not record a final selection yet. The required side-by-side evaluation against oapi-codegen must be completed before PR 3 is considered done.

The comparison input is the committed OpenAPI contract at `openapi/openapi.yaml`, generated from the TypeSpec source in `tsp/`. The contract currently describes 21 target `/api/v1` operations. Both generators must consume the same committed document and use pinned versions.

## Evidence collected so far

### ogen

- Prototype version: `ogen v1.25.0`.
- The full contract contains recursive expression schemas. In the first full-contract baseline run, unmodified ogen v1.25.0 stopped earlier with a generated schema-name collision (`PoolTermExprSpecSum`) while processing `/api/v1/pools/executions/spawns`; it did **not** reach a conclusive full-contract recursion result. The separate scoped recursion experiments remain valid evidence for their tested Pool.Term shapes, but must not be conflated with the full-contract run.
- Isolated experiments against the actual `Pool.Term` request schemas demonstrated two possible recursion-handling changes:
  - skip optional fields when detecting value-recursion, because their generated Go representation is indirect;
  - add pointer indirection for required recursive fields, preserving the OpenAPI `required` contract.
- The required and optional experiments use the same scoped OpenAPI operations and parity tests for missing, present, explicit-null, and nested recursive values. Both passed their isolated CI workflows.
- The required-recursion behavior is tracked upstream in [ogen-go/ogen#1792](https://github.com/ogen-go/ogen/issues/1792); optional-edge handling is tracked in [#1790](https://github.com/ogen-go/ogen/issues/1790). Required nullable decoding is separately tracked in [#1791](https://github.com/ogen-go/ogen/issues/1791).
- Detailed reproducible experiments are in [go-engine#79](https://github.com/orglang/go-engine/pull/79). That PR is intentionally draft-only and is not an implementation PR.

**Current assessment:** ogen can represent the tested recursive contract with narrowly scoped generator changes, while preserving wire-level requiredness in the tested cases. This is promising, but it creates a dependency on upstream fixes or a maintained local patch until those fixes are released. The experiments do not prove correctness for every recursive OpenAPI shape.

### oapi-codegen

- Candidate: `oapi-codegen v2.4.1` (to be pinned in the reproducible prototype).
- Generation against the same committed OpenAPI document has not yet been run as part of this comparison.
- Generated client compilation, recursive schema representation, nullable/optional behavior, and runtime JSON contract tests remain to be measured.

**Current assessment:** no conclusion yet. Do not treat the absence of a recorded failure as evidence that oapi-codegen passes.

## Reproducible generation baseline

The workflow `.github/workflows/api-client-generator-comparison.yaml` runs both pinned generators against the same committed `openapi/openapi.yaml` and records that input's SHA-256 checksum. The first full-contract run (artifact from run [#43591782147](https://github.com/orglang/go-engine/actions/runs/43591782147)) recorded ogen v1.25.0 failing with a schema-name collision (`PoolTermExprSpecSum`) and oapi-codegen v2.4.1 generating output with an OpenAPI 3.1 support warning. The committed config now pins oapi-codegen v2.8.2, which adds initial OpenAPI 3.1 support; its result is pending. The workflow preserves logs and generated outputs as an artifact.

This first automation stage is intentionally a **generation baseline**, not a completed client comparison: it does not yet establish compile/runtime parity, integration ergonomics, or a final choice. Those checks remain required below. A successful generation job only means both generator processes ran and artifacts were collected; it does not imply the generated clients compile or behave correctly.

## Comparison matrix

| Criterion from workspace issue #114 | ogen v1.25.0 | oapi-codegen |
| --- | --- | --- |
| Same real OpenAPI input | Full-contract baseline attempted; stopped at `PoolTermExprSpecSum` name collision. Scoped recursion experiments separately passed | v2.4.1 generated; warned about unsupported OpenAPI 3.1. v2.8.2 run pending |
| Request/response schemas and operations | Full-contract generation blocked by schema-name collision; recursive behavior still requires isolated evaluation | v2.4.1 generated client; full v2.8.2 evaluation pending |
| Optional vs nullable fields | Tested for the four Pool.Term continuation fields in isolated experiments; required nullable behavior has a separate control experiment | Not yet run |
| Enums and oneOf/discriminators | Present in the contract; full-contract output still needs review after recursion handling | Not yet run |
| Auth/security schemes and parameters | Must be inventoried against the actual contract and verified in generated client API | Not yet run |
| Default/error responses | Must be inventoried and verified in generated client API | Not yet run |
| Custom net/http transport, context/cancellation | Not yet evaluated in a representative consumer | Not yet run |
| Typed responses/errors and extension points | Not yet evaluated | Not yet run |
| Package/import shape and consumer ergonomics | Not yet evaluated | Not yet run |
| Generated code size/readability and dependencies | Not yet measured | Not yet measured |
| Clean regeneration and drift stability | Isolated generation workflows are reproducible; maintained pipeline not yet established | Not yet run |

## Required work to finish PR 3

- [x] Pin generator versions and run a first full-contract generation baseline against the exact same `openapi/openapi.yaml` (ogen v1.25.0; oapi-codegen v2.4.1). The current OpenAPI 3.1-capable oapi-codegen v2.8.2 run is pending.
- [ ] Keep generated outputs and configuration isolated; do not migrate consumers.
- [ ] Record generation success/failure and any required workarounds for each generator.
- [ ] Compile both generated clients and add representative request/response tests using deterministic local fixtures.
- [ ] Exercise the contract constructs that are actually present: recursive schemas, optional/nullable properties, enums, unions/discriminators, parameters, security (if declared), and response/error shapes.
- [ ] Build a small usage example against a representative existing consumer to compare transport injection, context cancellation, typed results/errors, and extension points.
- [ ] Record generated-code size and dependency footprint.
- [ ] Make a final selection with explicit trade-offs. If ogen is selected, document the upstream issue dependencies and the temporary policy for pinned patches.

## Provisional direction

ogen remains the leading candidate because its generator model and the successful isolated recursion experiments fit the contract. This is **not yet the final decision**: the comparison is incomplete until oapi-codegen is generated and tested against the identical OpenAPI input and the client integration criteria are evaluated.

## Scope

This is generator evaluation only. It does not add the maintained generated client to `orglang/go-sdk`, change server adapters, or migrate consumers. Those belong to PRs 4 and 7 in workspace issue #114.
