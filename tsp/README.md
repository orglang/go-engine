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

- `OrgLang.Engine.Type` — type references and type-link specifications
- `OrgLang.Engine.Term` — term references, variable specifications, and records
- `OrgLang.Engine.Comp` — computation references, variable specifications, and snapshots
- `OrgLang.Engine.Proc` — process declarations, types, steps, and expression specifications
- `OrgLang.Engine.Pool` — pool declarations, executions, steps, and expression specifications

Service operations use concise qualified names such as `Proc.DeclSnap`,
`Pool.ExecSpec`, `Type.Ref`, `Term.Ref`, and `Comp.Ref`. Models are
placed in the namespace for the concept they describe, rather than a generic
`Common` namespace.

The four expression specifications remain distinct and recursive:
`Proc.TypeExprSpec`, `Proc.TermExprSpec`, `Pool.TypeExprSpec`, and
`Pool.TermExprSpec`. Their schemas correspond to the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages.

## Naming

The Go wire prefix `comp` denotes computation, not component. The TypeSpec
namespace is therefore `Comp`, with concise model names such as `Comp.Ref`,
`Comp.VarSpec`, and `Comp.Snap`. Wire property names such as `comp_id`
and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
