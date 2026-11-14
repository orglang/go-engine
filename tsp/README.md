# Engine HTTP API contract

This directory is the source of truth for the HTTP API served by `go-engine`.
The OpenAPI document is generated into `openapi/openapi.yaml`; do not edit
that generated file by hand. The generation workflow also refreshes
`package-lock.json` from `package.json`.

## Generate

Requires Node.js 22 and npm. After the generation workflow has created or refreshed `package-lock.json`, install the exact dependency versions and generate OpenAPI:

```sh
cd tsp
npm ci
cd ..
task specs:generate
```

The CI workflow installs dependencies once per job before running `task specs:check`.

The contract describes 21 target `/api/v1` operations. Its resource paths are
intentionally normalized and are not yet the paths registered by the Echo server.
The server routes and existing clients remain unchanged until a separate migration.
The `/ssr` web presenter routes are intentionally out of scope.

## Generate the Echo HTTP server

The Go HTTP server boundary is generated from the emitted OpenAPI document with `oapi-codegen` and its Echo v4 integration. The generator version is pinned in `gear.mod`; the configuration and `go:generate` directive live in `internal/app/web/api`.

After generating OpenAPI, run:

```sh
task api:generate
```

The generated `ServerInterface` and route registration are transport-layer artifacts. Implementations belong in Echo adapters; generated files must not contain handwritten business logic. This step does not migrate the currently registered routes or claim that the normalized contract paths are already served.

## Model organization

TypeSpec namespaces use `OrgLang` as the root and follow domain concepts:

- `OrgLang.Any` — shared semantic identities (`TypeRef`, `TermRef`, `CompRef`) and the shared `ExecSpec`
- `OrgLang.Any.Type` — shared type-link specification
- `OrgLang.Any.Conn` — shared connection variable specs and records, defined in `any/conn/var.tsp`
- `OrgLang.Any.Exec.VarSpec` — execution variable specification shared by execution requests
- `OrgLang.Proc` — process declarations and snapshots
- `OrgLang.Proc.Type` / `OrgLang.Proc.Term` — process type and term expressions; type definitions live in `Proc.Type.DefSpec` / `DefSnap`
- `OrgLang.Pool` — pool declarations and snapshots, computation snapshots and execution steps
- `OrgLang.Pool.Type` / `OrgLang.Pool.Term` — pool type and term expressions; type definitions live in `Pool.Type.DefSpec` / `DefSnap`

The service is named `OrgLang Go Engine API`; `Engine` is not repeated as a namespace segment.
Declaration models use `DecSpec`/`DecSnap`, following the existing Go
implementation's `termdec` packages and types. Type definitions use
`Type.DefSpec`/`DefSnap`. Process and pool computation models remain
grouped in `Proc.Comp` and `Pool.Comp`.

## Expression polymorphism

The four expression specifications are discriminated unions:
`Proc.Type.ExprSpec`, `Proc.Term.ExprSpec`, `Pool.Type.ExprSpec`, and
`Pool.Term.ExprSpec`. Each union has one model per variant and serializes
the variant discriminator inline as the `kind` property, without an extra
payload envelope. A request must contain exactly one variant; the `kind`
value identifies that variant. Variant models are deliberately distinct even
when their payload fields currently have the same shape.

These models reflect the concepts handled by the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages, but are not
claimed to be wire-compatible with the current SDK DTOs. Several JSON names
are intentionally normalized in this target contract (for example
`chnl_ph` → `conn_ph`, `exp_vk` → `expr_ck`, `cont_es` → `cont_expr`,
and `val_es` → `val_expr`). The server and SDK migration must account for
these wire-level changes.

## Naming

The Go wire prefix `comp` denotes computation, not component. Shared
computation identity is `Any.CompRef`; context-specific computation models
remain grouped in their owning `Proc.Comp` and `Pool.Comp` namespaces.
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
any part conflicts or violates expression constraints, the whole operation
is rejected. A successful extension creates a new term revision.

Type definitions remain separate resources, retrieved through
`GET /api/v1/processes/types/{id}` and `GET /api/v1/pools/types/{id}`.
There are no type-definition `/definition` subresources or PUT operations
in this contract.

## Response conventions

POST operations that create a resource return `201 Created` with the created
resource's reference (`TypeRef`, `TermRef`, or `CompRef`) rather than its
full snapshot. Program submission is command-style and returns `201` without
a body. PATCH extension operations return `200 OK` with the new reference. Read operations return the requested snapshot, expression, or
reference collection. Command-style operations that do not create a distinct
resource return an empty `200 OK` or `204 No Content` response as specified
by the operation.

## Revision query

Single-resource retrieval operations that expose revisioned data accept an
optional `rn` query parameter. When omitted, the contract describes retrieval
of the latest available revision. When supplied, it identifies the requested
revision. Historical revision persistence is not asserted by this
specification.
