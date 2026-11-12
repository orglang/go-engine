# Engine HTTP API contract

This directory is the source of truth for the HTTP API served by `go-engine`.
The OpenAPI document is generated into `openapi/openapi.yaml`; do not edit
that generated file by hand.

## Generate

Requires Node.js and npm:

```sh
cd tsp
npm install
npm run build
```

The contract currently describes the 13 registered `/api/v1` operations.
The `/ssr` web presenter routes are intentionally out of scope.

## Model organization

TypeSpec namespaces use `OrgLang` as the root and follow domain concepts:

- `OrgLang.Any` — shared semantic identities (`TypeRef`, `TermRef`, `CompRef`) and the shared `ExecSpec`
- `OrgLang.Any.Type` — shared type-link specification
- `OrgLang.Any.Conn` — shared connection variable specs and records
- `OrgLang.Any.Exec.VarSpec` — execution variable specification shared by execution requests
- `OrgLang.Proc` — process declarations
- `OrgLang.Proc.Type` / `OrgLang.Proc.Term` — process type and term expressions; type definitions live in `Proc.Type.DefSpec` / `DefSnap`
- `OrgLang.Pool` — pool declarations, computation snapshots and execution steps
- `OrgLang.Pool.Type` / `OrgLang.Pool.Term` — pool type and term expressions; type definitions live in `Pool.Type.DefSpec` / `DefSnap`

The service is named `OrgLang Go Engine API`; `Engine` is not repeated as a namespace segment.
Declaration models use `DecSpec`/`DecSnap`, following the existing Go
implementation's `termdec` packages and types. Type definitions use
`Type.DefSpec`/`Type.DefSnap`. Process and pool computation models are
grouped in `Proc.Comp` and `Pool.Comp`.

The four expression specifications are distinct and recursive:
`Proc.Type.ExprSpec`, `Proc.Term.ExprSpec`, `Pool.Type.ExprSpec`, and
`Pool.Term.ExprSpec`. Their schemas correspond to the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages.

## Naming

The Go wire prefix `comp` denotes computation, not component. Shared
computation identity is `Any.CompRef`; context-specific computation models
are grouped in their owning `Proc.Comp` and `Pool.Comp` files. Wire property
names such as `comp_id` and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
