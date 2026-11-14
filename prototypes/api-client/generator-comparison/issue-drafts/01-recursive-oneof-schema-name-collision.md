## Description

ogen v1.25.0 fails to generate a recursive discriminated `oneOf` schema when an object schema is itself a `oneOf` and one of its variants refers back to that outer schema.

The failure is a generated type-name collision:

```text
schema name conflict: "ExprSum"
```

This is distinct from the infinite-recursion failure for required recursive fields. The reproduction below makes the recursive property optional, so the generated Go type graph can be finite through pointer indirection. This isolates the `ExprSum` collision from required-field recursion.

## OpenAPI specification

The following OpenAPI 3.1 document reproduces the problem:

```yaml
openapi: 3.1.0
info:
  title: Recursive oneOf schema name collision
  version: 1.0.0
paths:
  /expressions:
    post:
      operationId: createExpression
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Expr'
      responses:
        '204':
          description: No content
components:
  schemas:
    Expr:
      type: object
      oneOf:
        - $ref: '#/components/schemas/LeafExpr'
        - $ref: '#/components/schemas/BranchExpr'
      discriminator:
        propertyName: kind
        mapping:
          leaf: '#/components/schemas/LeafExpr'
          branch: '#/components/schemas/BranchExpr'
    LeafExpr:
      type: object
      required:
        - kind
      properties:
        kind:
          type: string
          const: leaf
    BranchExpr:
      type: object
      required:
        - kind
      properties:
        kind:
          type: string
          const: branch
        next:
          $ref: '#/components/schemas/Expr'
```

The recursive schema graph is:

```text
Expr --oneOf--> BranchExpr --optional--> Expr
```

The `next` property is intentionally optional. Making it required exercises a separate limitation covered by [ogen-go/ogen#1792](https://github.com/ogen-go/ogen/issues/1792).

## Steps to reproduce

Save the document as `recursive-oneof.yaml` and run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-recursive-oneof \
  --package probe \
  --clean recursive-oneof.yaml
```

## Expected behavior

Code generation succeeds. The outer `Expr` schema and its generated sum representation are registered consistently, and the recursive reference from `BranchExpr` resolves to the already-known `Expr` type.

## Actual behavior

ogen v1.25.0 fails during generation with:

```text
schema name conflict: "ExprSum"
```

The failure is reproducible with the optional recursive property, without relying on required-recursion handling.

## Suspected cause

The generator processes `oneOf`/ `anyOf` composition before registering the outer object schema. While collecting the sum variants, it recursively generates `BranchExpr`. Its `next` property refers back to `Expr`, so generation re-enters the outer composition and attempts to register another `ExprSum` type.

In other words, the outer schema is not available in the reference registry early enough to break the recursive generation path. The repeated sum registration then fails when the generated types are stored.

This is a source-level diagnosis based on ogen v1.25.0. It should be treated as a hypothesis until confirmed by an upstream fix or a focused regression test.

## Proposed change

Register the outer object schema before traversing its `oneOf`/`anyOf` variants. When a variant refers back to the outer schema, reference resolution should reuse that pre-registered object instead of entering composition generation again and registering a duplicate sum type.

The registration should preserve the schema's object representation (struct or map, depending on its properties and additional-properties configuration). The normal object-generation path should then reuse the pre-registered IR type rather than registering it a second time.

A regression test should verify that:

1. The optional-recursion reproduction generates successfully.
2. The generated package compiles.
3. The recursive reference resolves to the outer schema's object type, without a duplicate `ExprSum`.
4. Existing non-recursive discriminated `oneOf` schemas continue to generate as before.

## Experimental validation

A local candidate change in ogen v1.25.0 pre-registers object schemas that contain `oneOf`/`anyOf` before generating their variants.

The isolated CI experiment confirmed both sides of the reproduction:

- Without the candidate change, the minimal document fails with `schema name conflict: "ExprSum"`.
- With the candidate change, ogen generates the client and `go test ./client/...` compiles the generated package.

This validates the collision workaround for the optional-recursion fixture. It is not yet an upstream fix, and the experiment does not establish runtime JSON behavior or compatibility across a broader schema set.

## Related issues

- [ogen-go/ogen#1792](https://github.com/ogen-go/ogen/issues/1792) concerns required recursive fields and the need for pointer indirection while preserving JSON requiredness. That is a separate failure mode; the reproduction here uses an optional recursive edge specifically to avoid it.
- [orglang/go-engine#79](https://github.com/orglang/go-engine/pull/79) contains a related prototype and issue-draft format for recursion behavior.
