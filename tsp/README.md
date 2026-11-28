# Engine HTTP API contract

This directory is the source of truth for the HTTP API served by `go-engine`.
The OpenAPI document is generated into `openapi/openapi.yaml`; do not edit
that generated file by hand. The generation workflow also refreshes
`package-lock.json` from `package.json`.

## Generate

Requires Node.js 26.10.0. The repository's mise configuration installs Node.js,
Go, Task, and the Go code-generation tools:

```sh
mise install
mise run deps:install
task specs:generate
```

The `deps:install` task runs `npm ci` in this directory alongside downloading
Go modules. It installs the exact TypeSpec dependency versions from
`package-lock.json`.

The CI workflow installs dependencies once per job before running
`task specs:check`.

## Generate the Echo HTTP server

The Go HTTP server boundary is generated from the emitted OpenAPI document with
`oapi-codegen` and its Echo v5 integration. The generator version is pinned in
`mise.toml`; its configuration lives in `internal/app/web/api/cfg.yaml`.

The application and all existing HTTP adapters use Echo v5. The server
configuration enables `echo5-server`, ensuring the generated `ServerInterface`,
route registration, and production runtime share the same Echo API. The E2E
client configuration remains framework-independent and continues to generate
only the HTTP client and models.

The `task sources:generate` task generates Goverter adapters, the Echo server
boundary, and the E2E HTTP client from `openapi/openapi.yaml`. The E2E client
uses a separate `oapi-codegen` configuration and package under
`internal/lib/e2e/api`; its generated models are intentionally isolated from
the server package.

The generated `ServerInterface` and route registration are transport-layer
artifacts. Implementations belong in Echo adapters; generated files must not
contain handwritten business logic. Echo adapters implement the contract
operations and convert generated transport DTOs to engine-owned application
and domain types.

## Model organization

TypeSpec namespaces use `OrgLang` as the root and follow domain concepts:

- `OrgLang.Any` — shared semantic identities (`TypeRef`, `TermRef`, `RootRef`) and the shared `EvalSpec`
- `OrgLang.Any.Type` — shared type-link specification
- `OrgLang.Any.Conn` — shared connection variable specs and records, defined in `any/conn/var.tsp`
- `OrgLang.Any.Eval.VarSpec` — evaluation variable specification used when creating a pool execution
- `OrgLang.Proc` — process declarations and snapshots
- `OrgLang.Proc.Type` / `OrgLang.Proc.Term` — process type and term expressions; type definitions live in `Proc.Type.DefSpec` / `DefSnap`
- `OrgLang.Pool` — pool declarations and snapshots, execution snapshots and execution steps
- `OrgLang.Pool.Type` / `OrgLang.Pool.Term` — pool type and term expressions; type definitions live in `Pool.Type.DefSpec` / `DefSnap`

The service is named `OrgLang Go Engine API`; `Engine` is not repeated as a namespace segment.
Declaration models use `DecSpec`/`DecSnap`, following the existing Go
implementation's `termdec` packages and types. Type definitions use
`Type.DefSpec`/`DefSnap`. Process and pool execution models remain
grouped in `Proc.Exec` and `Pool.Exec` namespaces.

## Expression polymorphism

The four expression specifications are discriminated unions:
`Proc.Type.ExprSpec`, `Proc.Term.ExprSpec`, `Pool.Type.ExprSpec`, and
`Pool.Term.ExprSpec`. Each union has one model per variant and serializes
the variant discriminator inline as the `kind` property, without an extra
payload envelope. A request must contain exactly one variant; the `kind`
value identifies that variant. Variant models are deliberately distinct even
when their payload fields currently have the same shape.

These models describe the engine's process and pool expression variants
independently of any client implementation. Several JSON names are
intentionally normalized in the contract (for example `chnl_ph` → `conn_ph`,
`exp_vk` → `expr_ck`, `cont_es` → `cont_expr`, and `val_es` → `val_expr`).
HTTP adapters map these generated transport models to engine-owned types.

## Naming

The Go wire prefix `comp` denotes legacy computation prefix, not component.
Shared root identity is `Any.RootRef`; context-specific computation
models remain grouped in their owning `Proc.Exec` and `Pool.Exec` namespaces.
The HTTP resource paths use the common term `executions` for both processes
and pools, while retaining the domain model name `Comp`. Nested collection
routes return arrays of references; nested `declaration` routes return one
declaration snapshot.

## Definition retrieval and extension

Type definition snapshots remain available at
`GET /api/v1/processes/types/{id}` and
`GET /api/v1/pools/types/{id}`.

The term definition subresource is consistent for processes and pools:

- `GET /api/v1/processes/{id}/definition` returns `Proc.Term.ExprSpec`.
- `PATCH /api/v1/processes/{id}/definition` accepts a `Proc.Term.ExprSpec`
  fragment and returns the updated `Any.TermRef`.
- `GET /api/v1/pools/{id}/definition` returns `Pool.Term.ExprSpec`.
- `PATCH /api/v1/pools/{id}/definition` accepts a `Pool.Term.ExprSpec`
  fragment and returns the updated `Any.TermRef`.

Both GET operations accept the optional `rn` revision query. Both PATCH
operations apply the supplied expression as an additive extension to the
existing term definition; they do not replace it. Existing expression parts
cannot be removed or changed. The extension must be applied atomically: if
any part conflicts or violates expression constraints, the whole operation is
rejected. A successful extension advances the owning Description's `DescRN`
exactly once; it does not create a separate Term revision.

Type definitions remain separate resources, retrieved through
`GET /api/v1/processes/types/{id}` and `GET /api/v1/pools/types/{id}`.
There are no type-definition `/definition` subresources or PUT operations
in this contract.

## Response conventions

POST operations that create a resource return `201 Created` with the created
resource's reference (`TypeRef`, `TermRef`, or `RootRef`) rather than its
full snapshot. Program submission is command-style and returns `201` without
a body. PATCH extension operations return `200 OK` with the new reference.
Read operations return the requested snapshot, expression, or reference collection.
Command-style operations that do not create a distinct resource return an
empty `200 OK` or `204 No Content` response as specified by the operation.

## Revision query

Single-resource retrieval operations that expose revisioned data accept an
optional `rn` query parameter. For Term Definition retrieval, `rn` is the
owning Description's `DescRN`. When omitted, the current expression is
returned. Historical Description snapshots are not persisted, so a supplied
`rn` can identify only the current revision.
