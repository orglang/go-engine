# Go client generator comparison

Tracks PR 3 from [orglang/go-workspace#114](https://github.com/orglang/go-workspace/issues/114).

## Status

**Draft / evaluation in progress.** ogen is the leading candidate, but this document does not record a final selection yet. The required side-by-side evaluation against oapi-codegen must be completed before PR 3 is considered done.

The comparison input is the committed OpenAPI contract at `openapi/openapi.yaml`, generated from the TypeSpec source in `tsp/`. The contract currently describes 21 target `/api/v1` operations. Both generators must consume the same committed document and use pinned versions.

## Evidence collected so far

### Latest runtime integration evidence — comparison run #37

- Commit: [0c07b5d](https://github.com/orglang/go-engine/commit/0c07b5d6331fe84fb4d531e0606416071944d61f); [CI run #43596829412](https://github.com/orglang/go-engine/actions/runs/43596829412) completed successfully. All 12 main workflow steps passed, including the ogen diagnostic generation and package test, full-contract oapi-codegen generation, generated-package tests, and artifact upload.
- The oapi-codegen test suite now verifies a custom `http.Client` transport via `WithHTTPClient`, request editing via `WithRequestEditorFn` (including an injected header), and parsing a typed `201` response through `NewClientWithResponses`.
- A separate test confirms that cancellation of the request context propagates through the custom transport as `context.Canceled`.
- Together with run #36, the tested client behaviors now include JSON and text/plain request bodies, path parameter escaping, typed `201` success, empty `204`, raw body/status for an unmodeled `422`, custom transport, request editor, and context cancellation. This remains focused test coverage, not a complete consumer integration or a comprehensive error contract.
- The full-contract ogen baseline still fails with `schema name conflict: "PoolTermExpressionDefinitionSum"`; the green workflow records this expected failure and does not imply successful full-contract ogen generation. The filtered two-operation ogen diagnostic still generates and compiles.
- Run #37 artifact sizes: full oapi-codegen output directory 224,588 bytes (generated `client.gen.go`: 206,559 bytes); filtered ogen diagnostic output directory 98,573 bytes. These remain incomparable because the generated scopes differ.

### ogen

- Prototype version: `ogen v1.25.0`.
- The full contract contains recursive expression schemas. In the full-contract baseline, unmodified ogen v1.25.0 stops with a generated schema-name collision (initially `PoolTermExprSpecSum`) while processing `/api/v1/pools/executions/spawns`; it does **not** reach a conclusive full-contract recursion result.
- A separate diagnostic now filters the OpenAPI document to `/api/v1/pools/executions/spawns` and `/api/v1/pools/executions/steps`, removes the four `cont_expr` edges from the Pool.Term variants, and runs ogen v1.25.0. Generation succeeds and `go test ./ogen-no-recursion/...` passes, confirming that this generated package compiles. CI evidence: [comparison run #43593945399](https://github.com/orglang/go-engine/actions/runs/43593945399), including the successful “Diagnose recursive Pool.Term.ExprSpec edges” and “Test generated ogen package” steps. This is a deliberately narrowed diagnostic: it does not establish full-contract generation, and it does not test runtime behavior of the generated package.
- Isolated experiments against the actual `Pool.Term` request schemas demonstrated two possible recursion-handling changes:
  - skip optional fields when detecting value-recursion, because their generated Go representation is indirect;
  - add pointer indirection for required recursive fields, preserving the OpenAPI `required` contract.
- The required and optional experiments use the same scoped OpenAPI operations and parity tests for missing, present, explicit-null, and nested recursive values. Both passed their isolated CI workflows.
- The required-recursion behavior is tracked upstream in [ogen-go/ogen#1792](https://github.com/ogen-go/ogen/issues/1792); optional-edge handling is tracked in [#1790](https://github.com/ogen-go/ogen/issues/1790). Required nullable decoding is separately tracked in [#1791](https://github.com/ogen-go/ogen/issues/1791).
- Detailed reproducible experiments are in [go-engine#79](https://github.com/orglang/go-engine/pull/79). That PR is intentionally draft-only and is not an implementation PR.

**Current assessment:** ogen can represent the tested recursive contract with narrowly scoped generator changes, while preserving wire-level requiredness in the tested cases. This is promising, but it creates a dependency on upstream fixes or a maintained local patch until those fixes are released. The experiments do not prove correctness for every recursive OpenAPI shape.

### oapi-codegen

- Candidate: `oapi-codegen v2.8.2` (pinned in the reproducible workflow).
- The v2.8.2 full-contract generation completed successfully against the same committed OpenAPI input. In comparison run [#43595313547](https://github.com/orglang/go-engine/actions/runs/43595313547), the full generated package passed `go test ./oapi-codegen/...` with deterministic JSON contract tests for the recursive Hire-to-Fire continuation and missing/explicit-null `cont_expr`. Follow-up run [#43595575239](https://github.com/orglang/go-engine/actions/runs/43595575239) added passing round-trip tests for two additional discriminated unions: `Pool.Type.ExprSpec` (`link`) and `Proc.Term.ExprSpec` (`close`). Run [#43596153549](https://github.com/orglang/go-engine/actions/runs/43596153549) additionally verifies request query serialization for `EngineApiGetPoolTypeDefinition`: a provided `Rn` emits `rn=42`, while nil omits the parameter. The workflow's full generated-package test step passed. The latter two behaviors show that ordinary generated model decoding does not enforce OpenAPI required-property presence or non-nullability. This is evidence for these specific model behaviors, not comprehensive runtime parity or validation of all contract constraints.
- The same run recorded generated directory sizes: 206,559 bytes for the full oapi-codegen output and 98,573 bytes for the scoped ogen diagnostic output. These figures are **not a valid size comparison** because the generated scopes differ (full contract vs. two operations with four recursion edges removed).

**Current assessment:** full-contract generation and package compilation succeed, and focused JSON tests cover recursive union round-trip plus the decoder's behavior for both missing and explicit-null `cont_expr`. The generated model decoder does not enforce OpenAPI required-property presence or non-nullability for these cases; callers needing contract validation must provide a separate validation layer. Recursive union round-trip is demonstrated only for the tested Hire-to-Fire example.

## Reproducible generation baseline

The workflow `.github/workflows/api-client-generator-comparison.yaml` runs both pinned generators against the same committed `openapi/openapi.yaml` and records that input's SHA-256 checksum. It also runs a separate scoped ogen diagnostic: after filtering to the two Pool execution operations and removing the four `cont_expr` recursion edges, it requires both generation and Go package compilation to pass. The first full-contract run (artifact from run [#43591782147](https://github.com/orglang/go-engine/actions/runs/43591782147)) recorded ogen v1.25.0 failing with a schema-name collision (`PoolTermExprSpecSum`) and oapi-codegen v2.4.1 generating output with an OpenAPI 3.1 support warning. The follow-up run [#43591852634](https://github.com/orglang/go-engine/actions/runs/43591852634) completed with oapi-codegen v2.8.2 generation successful; ogen still reported the original name collision. The workflow preserves logs and generated outputs as an artifact.

This first automation stage is intentionally a **generation baseline**, not a completed client comparison: it does not yet establish compile/runtime parity, integration ergonomics, or a final choice. Those checks remain required below. A successful generation job only means both generator processes ran and artifacts were collected; it does not imply the generated clients compile or behave correctly.

## ogen schema-name override experiment

The first ogen baseline failed while generating `/api/v1/pools/executions/spawns` because the generated sum type name `PoolTermExprSpecSum` collided. The first override changed the collision to `PoolTermExpressionSum`, proving that ogen v1.25.0 recognizes `x-ogen-name`, but not resolving it.

The follow-up experiment gives the component schema the name `PoolTermExpressionDefinition` and adds a distinct `x-ogen-name: PoolStepTermExpression` sibling to the `$ref` on `Pool.StepSpec.pool_expr`. This tests whether the conflicting generated sum types originate from the component schema and its use through the step property.

These extensions are generator metadata only; they do not change the JSON schema or wire contract, and oapi-codegen is expected to ignore them. They are applied to the committed OpenAPI document only for this experiment. The TypeSpec source has not yet been taught to emit them, so normal TypeSpec regeneration may remove them. Do not treat this as a durable source-of-truth fix until the generated Go names are verified and the extension is represented in the TypeSpec pipeline.

The follow-up workflow run [#43592376123](https://github.com/orglang/go-engine/actions/runs/43592376123) completed successfully as a workflow, but the saved `ogen.log` shows that ogen v1.25.0 still exits with `schema name conflict: "PoolTermExpressionDefinitionSum"` while generating the `/api/v1/pools/executions/spawns` request body. The component-level rename changed the conflicting generated name again; the property-level `$ref` sibling did not eliminate the collision. The green workflow is expected because the script captures the ogen exit status and uploads diagnostics rather than failing the job on this known generator failure. oapi-codegen v2.8.2 still completed generation successfully.

## Comparison matrix

| Criterion from workspace issue #114 | ogen v1.25.0 | oapi-codegen |
| --- | --- | --- |
| Same real OpenAPI input | Full-contract baseline still stops at a sum-type name collision (`PoolTermExpressionDefinitionSum` after schema-name overrides); scoped diagnostic package compiles after recursion-edge removal | v2.8.2 full-contract generation and package compilation passed; runtime checks pending |
| Request/response schemas and operations | Full-contract generation blocked by persistent schema-name collision; scoped two-operation package compiles after removing four `cont_expr` recursion edges | v2.8.2 full-contract package compiles; focused recursive JSON round-trip passes |
| Optional vs nullable fields | Tested for the four Pool.Term continuation fields in isolated experiments; required nullable behavior has a separate control experiment | Missing required `cont_expr` is accepted by model decoding in the focused test; broader optional/nullable matrix pending |
| Enums and oneOf/discriminators | Present; full-contract output still blocked by the name collision | No enum schemas are declared in the current contract; three discriminator variants now have passing JSON round-trip coverage (`Hire.cont_expr`, `Pool.Type.ExprSpec`, `Proc.Term.ExprSpec`; runs #43595313547 and #43595575239) |
| Auth/security schemes and parameters | Must be inventoried against the actual contract and verified in generated client API | Optional `rn` query parameter serialization is tested (present and omitted); security and broader parameter coverage pending |
| Default/error responses | Must be inventoried and verified in generated client API | Not yet run |
| Custom net/http transport, context/cancellation | Full-contract output blocked; diagnostic scope only | oapi-codegen tests pass for custom `http.Client` transport, `WithRequestEditorFn`, and `context.Canceled` propagation (run #43596829412) |
| Typed responses/errors and extension points | Not yet evaluated | Not yet run |
| Package/import shape and consumer ergonomics | Not yet evaluated | Not yet run |
| Generated code size/readability and dependencies | 98,573 bytes for filtered diagnostic output only; not comparable | 224,588 bytes full output directory; generated file 206,559 bytes. Scope mismatch prevents a valid size comparison; dependency footprint still to measure |
| Clean regeneration and drift stability | Isolated generation workflows are reproducible; maintained pipeline not yet established | Not yet run |

## Required work to finish PR 3

- [x] Pin generator versions and run full-contract generation against the exact same `openapi/openapi.yaml` (ogen v1.25.0; oapi-codegen v2.8.2).
- [x] Keep generated outputs and configuration isolated; do not migrate consumers.
- [x] Record generation success/failure and the scoped ogen diagnostic workaround; full-contract ogen generation remains blocked by the name collision.
- [x] Compile the generated oapi-codegen full-contract package and the scoped ogen diagnostic package.
- [x] Add deterministic JSON contract tests for recursive Hire-to-Fire continuation round-trip, missing/explicit-null `cont_expr`, and two additional `oneOf` discriminator round-trips (`Pool.Type.ExprSpec` and `Proc.Term.ExprSpec`; CI runs [#43595313547](https://github.com/orglang/go-engine/actions/runs/43595313547) and [#43595575239](https://github.com/orglang/go-engine/actions/runs/43595575239)).
- [x] Verify optional `rn` query serialization for a provided value and nil omission (CI run [#43596153549](https://github.com/orglang/go-engine/actions/runs/43596153549)).
- [ ] Exercise the contract constructs that are actually present: recursive schemas, optional/nullable properties, enums, unions/discriminators, parameters, security (if declared), and response/error shapes.
- [x] Exercise oapi-codegen transport injection, request editor, context cancellation, typed success, no-content success, and unmodeled error response behavior (runs [#43596671894](https://github.com/orglang/go-engine/actions/runs/43596671894) and [#43596829412](https://github.com/orglang/go-engine/actions/runs/43596829412)).
- [ ] Build a small usage example against a representative existing consumer; current tests use a focused generated-client harness, not a production consumer.
- [x] Record generated output sizes, with an explicit warning that the full oapi-codegen output and filtered ogen diagnostic are not comparable (run [#43596829412](https://github.com/orglang/go-engine/actions/runs/43596829412)).
- [ ] Compare dependency footprint and generated-code readability on equivalent scopes.
- [ ] Make a final selection with explicit trade-offs. If ogen is selected, document the upstream issue dependencies and the temporary policy for pinned patches.

## Provisional direction

ogen remains the leading candidate because its generator model and the successful isolated recursion experiments fit the contract. This is **not yet the final decision**: the comparison is incomplete until oapi-codegen is generated and tested against the identical OpenAPI input and the client integration criteria are evaluated.

## Scope

This is generator evaluation only. It does not add the maintained generated client to `orglang/go-sdk`, change server adapters, or migrate consumers. Those belong to PRs 4 and 7 in workspace issue #114.
