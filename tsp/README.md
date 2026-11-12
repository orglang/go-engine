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

## Revision query

The three retrieval operations accept an optional `revision` query parameter.
When omitted, the contract describes retrieval of the latest available revision.
When supplied, it identifies the requested revision. This is an API contract
statement; persistence support for historical revisions is not asserted by
this specification.

## Model fidelity

The expression payload is currently represented as an open object with a
required `kind` discriminator. The SDK contains separate recursive expression
variants for process types, pool types, process steps, and pool steps. Those
variant schemas should be expanded in a follow-up contract refinement rather
than silently treating the current open shape as exhaustive.
