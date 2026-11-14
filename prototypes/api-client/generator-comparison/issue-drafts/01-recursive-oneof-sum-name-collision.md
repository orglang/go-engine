# Recursive `oneOf` object schemas produce an `ExprSum` name collision

## Description

ogen v1.25.0 fails to generate a recursive object schema that declares both
`type: object` and `oneOf`. The schema's `oneOf` variants include an object
that refers back to the outer schema through an **optional** property.

Instead of resolving the recursive reference to the already-declared outer
object type, ogen attempts to generate the `oneOf` sum type again. The second
registration uses the same generated name and fails with:

```text
schema name conflict: "ExprSum"
```

This is a schema/type-registration collision. It is distinct from the
infinite-recursion detection problem for required recursive fields.

## OpenAPI specification

The following OpenAPI 3.1 document is a minimal reproduction. The recursive
property `BranchExpr.next` is deliberately optional: it is not listed in the
schema's `required` array.

```yaml
openapi: 3.1.0
info:
  title: Recursive oneOf name collision
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
      required: [kind]
      properties:
        kind:
          type: string
          const: leaf
    BranchExpr:
      type: object
      required: [kind]
      properties:
        kind:
          type: string
          const: branch
        next:
          $ref: '#/components/schemas/Expr'
```

The recursive path is:

```text
Expr --oneOf--> BranchExpr --optional--> Expr
```

The optional property is important: this reproduction does not require
ogen to represent a required recursive field through pointer indirection.

## Steps to reproduce

Save the document as `recursive-oneof.yaml` and run:

```sh
go run github.com/ogen-go/ogen/cmd/ogen@v1.25.0 \
  --target /tmp/ogen-recursive-oneof \
  --package probe \
  --clean recursive-oneof.yaml
```

## Expected behavior

Code generation succeeds. The outer `Expr` schema is registered once, and
the recursive reference from `BranchExpr.next` resolves to that type. The
generator emits one `ExprSum` definition for the `oneOf` union.

## Actual behavior

ogen v1.25.0 fails during schema generation with:

```text
schema name conflict: "ExprSum"
```

This was the optional-edge form of the prototype reproduction used in the
successful candidate-patch run below. The committed fixture was later changed
to make `BranchExpr.next` required for a separate experiment concerning
[ogen-go/ogen#1792](https://github.com/ogen-go/ogen/issues/1792).

## Suspected cause

The schema generator processes `oneOf`/`anyOf` before registering the
outer object schema. While generating the sum variants, it follows the
recursive reference back to the outer schema and enters the same
`oneOf` generation path again.

The outer sum type is therefore generated and registered more than once
under the same name (`ExprSum`).

A likely fix is to register the outer object schema before traversing its
`oneOf`/`anyOf` variants, so a recursive reference can resolve to the
in-progress object type rather than recursively starting sum generation.
This is a source-level hypothesis; the exact registration strategy belongs
to the generator implementation.

## Experimental evidence

The behavior was reproduced with ogen v1.25.0 in the isolated prototype:

- [Workflow run #55](https://github.com/orglang/go-engine/actions/runs/43599686693)
  confirmed the expected `schema name conflict: "ExprSum"` failure for the
  minimal recursive `oneOf` schema.
- A candidate patch that pre-registers the outer object before traversing
  the union variants was then tested with the recursive property left
  optional. [Workflow run #67](https://github.com/orglang/go-engine/actions/runs/43601032094)
  completed successfully: generation with the patched ogen source succeeded,
  and `go test ./client/...` compiled the generated package.

The latter is a focused generation/compilation check. It does not establish
runtime JSON decoding or marshaling behavior, nor does it prove that the
candidate fix covers every recursive `oneOf`/ `anyOf` schema shape.

## Relation to required recursive fields

This report is separate from
[ogen-go/ogen#1792](https://github.com/ogen-go/ogen/issues/1792)
(`gen: use pointer indirection for required recursive fields`).

- **This report:** the recursive property is optional; ogen fails because
  the outer object/union is registered more than once as `ExprSum`.
- **#1792:** the recursive path consists of required properties; ogen fails
  with an `infinite recursion: ... is required` error because the generated
  Go representation needs indirection.

The two failures have different symptoms and require separate regression
tests. A fix for required-field recursion should not be treated as evidence
that the `ExprSum` name collision is fixed, or vice versa.

## Suggested regression test

Keep this optional-edge schema as a standalone fixture and verify that:

1. Unpatched ogen v1.25.0 fails with the `ExprSum` name collision.
2. The candidate fix generates the client successfully.
3. The generated package compiles.

Runtime serialization tests can be added separately if the generator change
also changes generated model or decoder behavior.
