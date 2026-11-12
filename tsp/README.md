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

TypeSpec namespaces and files follow the Go SDK package boundaries:

- `OrgLang.Engine.Common` — shared references and variable records
- `OrgLang.Engine.Proc` — process declarations, types, steps, and expression specs
- `OrgLang.Engine.Pool` — pool declarations, executions, steps, and expression specs

For example, service operations refer to `Proc.DeclSnap`,
`Proc.TypeSpec`, `Pool.ExecSpec`, and `Common.TermRef`. Within each
namespace, model names do not repeat the package prefix.

The four expression specifications remain distinct and recursive:
`Proc.TypeExprSpec`, `Proc.TermExprSpec`, `Pool.TypeExprSpec`, and
`Pool.TermExprSpec`. Their schemas correspond to the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages.

## Naming

The Go wire prefix `comp` denotes computation, not component. In TypeSpec,
`Common.ComputationRef`, `Common.CompVarSpec`, and
`Pool.ComputationSnap` make that distinction explicit without repeating
the long word in every model name. The `CompVarSpec` name is deliberately
short; its `comp` prefix means computation. Wire property names such as
`comp_id` and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
