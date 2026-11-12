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

Models are split under `models/` to follow the corresponding Go SDK package
boundaries. Process models live under `models/proc/`; pool models live under
`models/pool/`; shared references and variable records live in
`models/common.tsp`. `main.tsp` contains service metadata and HTTP routes.

The four expression specifications are distinct and recursive:
`ProcTypeExprSpec`, `ProcTermExprSpec`, `PoolTypeExprSpec`, and
`PoolTermExprSpec`. Their schemas mirror the respective
`go-sdk/proc/typeexp`, `go-sdk/proc/termexp`,
`go-sdk/pool/typeexp`, and `go-sdk/pool/termexp` packages. Do not merge
these into a generic expression model.

## Naming

The `comp` prefix in the Go wire fields denotes computation, not component.
The TypeSpec model names therefore use `Computation` (for example,
`ComputationRef`, `ComputationVarSpec`, and `ComputationSnap`).
Wire property names such as `comp_id` and `comp_ref` remain unchanged.

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.
