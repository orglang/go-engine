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

- `OrgLang.Engine.Type` — shared type references and type-link specifications
- `OrgLang.Engine.Term` — shared term references, variable specifications, and records
- `OrgLang.Engine.Comp` — shared computation identity (`Comp.Ref`)
- `OrgLang.Engine.Proc` — process declarations, types, and expression specifications
- `OrgLang.Engine.Proc.Comp` — process computation snapshots and step specifications
- `OrgLang.Engine.Pool` — pool declarations and expression specifications
- `OrgLang.Engine.Pool.Comp` — pool computation specs, snapshots, and step specifications

The `Comp.Ref` model is shared because both process and pool computations use
the same computation identity. Models describing computation behavior or
payloads are scoped to their owning context: `Proc.Comp.*` and `Pool.Comp.*`.
Service operations use qualified names such as `Proc.Comp.Snap`,
`Proc.Comp.StepSpec`, `Pool.Comp.ExecSpec`, `Pool.Comp.Snap`,
`Pool.Comp.StepSpec`, `Type.Ref`, and `Term.Ref`.

The four expression specifications remain distinct and recursive:
`Proc.TypeExprSpec`, `Proc.TermExprSpec`, `Pool.TypeExprSpec`, and
`Pool.TermExprSpec`. Their schemas correspond to the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages.

## Naming

The Go wire prefix `comp` denotes computation, not component. The TypeSpec
namespace is therefore `Comp`, with the shared identity model `Comp.Ref`.
Context-specific models live under `Proc.Comp` and `Pool.Comp`. Wire
property names such as `comp_id` and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
