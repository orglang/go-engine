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

TypeSpec namespaces follow the Go SDK package boundaries and shared concepts:

- `OrgLang.Engine.Any.Type`, `OrgLang.Engine.Any.Term`, and `OrgLang.Engine.Any.Comp` — shared type, term, and computation identities
- `OrgLang.Engine.Type` — type-link specifications
- `OrgLang.Engine.Term` — variable specifications and records
- `OrgLang.Engine.Proc` — process declarations, types, computation snapshots, step specs, and expression specifications
- `OrgLang.Engine.Pool` — pool declarations, computation snapshots/specs/steps, and expression specifications

Only semantic identities shared by process and pool contexts live under `Any`.
Context-specific computation models are direct members of their owning namespaces:
`Proc.CompSnap`, `Proc.StepSpec`, `Pool.CompSnap`, `Pool.VarSpec`,
`Pool.ExecSpec`, and `Pool.StepSpec`.

The four expression specifications remain distinct and recursive:
`Proc.TypeExprSpec`, `Proc.TermExprSpec`, `Pool.TypeExprSpec`, and
`Pool.TermExprSpec`. Their schemas correspond to the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages.

## Naming

The Go wire prefix `comp` denotes computation, not component. Shared
computation identity is `Any.Comp.Ref`; context-specific snapshots and step
specifications are named directly in `Proc` and `Pool`. Wire property names
such as `comp_id` and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
