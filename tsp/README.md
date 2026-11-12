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

TypeSpec namespaces use `OrgLang` as the root and follow the Go SDK package boundaries and shared concepts:

- `OrgLang.Any` (`TypeRef`, `TermRef`, and `CompRef`) — shared type, term, and computation identities
- `OrgLang.Conn` — connection variable specifications and records
- `OrgLang.Type` — type-link specifications
- `OrgLang.Proc` — process declarations, definitions, computation snapshots, and step specs
- `OrgLang.Proc.Type` / `OrgLang.Proc.Term` — process type and term expression specifications
- `OrgLang.Pool` — pool declarations, definitions, computation snapshots/specs, and steps
- `OrgLang.Pool.Type` / `OrgLang.Pool.Term` — pool type and term expression specifications

Only semantic identities shared by process and pool contexts live under `Any`.
The service is named `OrgLang Go Engine API`; `Engine` is not repeated as a namespace segment.
Context-specific computation models are direct members of their owning namespaces:
`Proc.CompSnap`, `Proc.StepSpec`, `Pool.CompSnap`, `Pool.VarSpec`,
`Pool.ExecSpec`, and `Pool.StepSpec`.

Declaration models use `DecSpec`/`DecSnap`, following the existing Go
implementation's `termdec` packages and `DecSpec`/ `DecSnap` types.
Type-definition models use `DefSpec`/`DefSnap`. Connection variable
specifications and records are shared under `Conn`.

The four expression specifications are distinct and recursive:
`Proc.Type.ExprSpec`, `Proc.Term.ExprSpec`, `Pool.Type.ExprSpec`, and
`Pool.Term.ExprSpec`. Their schemas correspond to the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages.

## Naming

The Go wire prefix `comp` denotes computation, not component. Shared
computation identity is `Any.CompRef`; context-specific snapshots and step
specifications are named directly in `Proc` and `Pool`. Wire property names
such as `comp_id` and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
